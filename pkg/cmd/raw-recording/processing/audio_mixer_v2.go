package processing

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

// segmentMixWarnThreshold bounds how many segments we expect to feed into a
// single amix. With 1-3 segments per track this is only reached by pathological
// calls (very many participants or frequent renegotiation); we log rather than
// chunk, since amix handles hundreds of inputs and chunking would reintroduce
// intermediate files.
const segmentMixWarnThreshold = 200

type AudioMixerV2 struct {
	logger *ProcessingLogger
}

func NewAudioMixerV2(logger *ProcessingLogger) *AudioMixerV2 {
	return &AudioMixerV2{logger: logger}
}

// MixAllAudioTracks decodes every matching audio segment and mixes them into a
// single output file. Each segment is fed directly into the mixer positioned at
// its absolute first-packet timestamp (adelay), so inter-segment gaps become
// silence without any per-track concatenation or gap-filler files.
func (p *AudioMixerV2) MixAllAudioTracks(config *AudioMixerConfig, metadata *RecordingMetadata) (*string, error) {
	p.overrideConfig(config)

	p.logger.Infof("Extracting all matching audio segments...")
	segments, err := p.extractSegmentsForMix(config, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to extract audio segments: %w", err)
	}

	hasWindow := config.RecordingEnd.After(config.RecordingStart)

	fileOffsets, firstPacketTime := p.offset(segments)
	if len(fileOffsets) == 0 {
		p.logger.Warnf("No audio tracks found; skipping composite output")
		return nil, nil
	}

	p.logger.Infof("Found %d audio segments to mix", len(fileOffsets))

	// Clean up the per-segment container files after mixing (optional).
	if config.WithCleanup {
		defer func(offsets *[]*FileOffset) {
			for _, fileOffset := range *offsets {
				p.logger.Infof("Cleaning up temporary file: %s", fileOffset.Name)
				if err := os.Remove(fileOffset.Name); err != nil {
					p.logger.Warnf("Failed to clean up temporary file %s: %v", fileOffset.Name, err)
				}
			}
		}(&fileOffsets)
	}

	// Anchor the mix to the recording window without silence files: recordingDurationSeconds
	// pads the tail to the window end (trailing silence), and leadingGapMs is the
	// silence before the first track (added to the mixer's offset rebasing).
	leadingGapMs, recordingDurationSeconds := int64(0), float64(0)
	if hasWindow {
		recordingDurationSeconds = config.RecordingEnd.Sub(config.RecordingStart).Seconds()
		if gap := time.UnixMilli(firstPacketTime).Sub(config.RecordingStart); gap > 0 {
			leadingGapMs = gap.Milliseconds()
		}
	}

	outputFile := p.buildFilename(config, metadata)

	if err := runFFmpegCommand(generateMixAudioFilesArguments(outputFile, config.Format, fileOffsets, leadingGapMs, recordingDurationSeconds), p.logger); err != nil {
		return nil, fmt.Errorf("failed to mix audio files: %w", err)
	}

	p.logger.Infof("Successfully created mixed audio file: %s", outputFile)

	return &outputFile, nil
}

func (p *AudioMixerV2) overrideConfig(config *AudioMixerConfig) {
	if !slices.Contains(supportedFormats[:], config.Format) {
		p.logger.Warnf("Audio format %s not supported, fallback to default %s", config.Format, DefaultFormat)
		config.Format = DefaultFormat
	}
}

// extractSegmentsForMix decodes every audio segment matching the mix filter into
// its own container file and returns those segments (each with its ContainerPath
// set) ready to be positioned and mixed.
//
// Unlike the per-track extraction path it neither concatenates segments nor
// synthesizes inter-segment gap fillers: the mixer's adelay positioning turns
// inter-segment gaps into silence on its own.
func (p *AudioMixerV2) extractSegmentsForMix(config *AudioMixerConfig, metadata *RecordingMetadata) ([]*SegmentInfo, error) {
	mediaType := ""
	if !config.WithScreenshare {
		mediaType = MediaTypeUser
	}

	tracks := FilterTracks(metadata.Tracks, "", "", "", trackKindAudio, mediaType)
	if len(tracks) == 0 {
		p.logger.Warnf("No audio tracks found matching the filter criteria")
		return nil, nil
	}

	var segments []*SegmentInfo
	for _, t := range tracks {
		segments = append(segments, t.Segments...)
	}

	if len(segments) > segmentMixWarnThreshold {
		p.logger.Warnf("Mixing %d audio segments in a single pass, above the expected threshold of %d", len(segments), segmentMixWarnThreshold)
	}

	// Resolve each segment's container path during a single directory walk, then
	// decode the matching rtpdumps into their per-segment container file.
	accept := func(path string, info os.FileInfo) (*SegmentInfo, bool) {
		for _, s := range segments {
			if s.ContainerPath != "" {
				continue // already matched to another walked file
			}
			if strings.Contains(info.Name(), s.metadata.BaseFilename) {
				extension, suffix := oggExtension, oggSuffix
				abs, _ := filepath.Abs(path)

				s.RtpDumpPath = abs
				s.SdpPath = strings.ReplaceAll(abs, suffixRtpDump, suffixSdp)
				s.ContainerExt = extension
				s.ContainerPath = strings.ReplaceAll(abs, suffixRtpDump, suffix)
				return s, true
			}
		}
		return nil, false
	}

	if err := ConvertDirectory(config.WorkDir, accept, false, p.logger); err != nil {
		return nil, fmt.Errorf("failed to convert audio segments: %w", err)
	}

	return segments, nil
}

// offset builds one FileOffset per converted segment, keyed by the segment's
// absolute first-packet timestamp, and returns the earliest such timestamp
// (firstPacketTime). Segments without a container path (conversion failed) are skipped.
func (p *AudioMixerV2) offset(segments []*SegmentInfo) (offsets []*FileOffset, firstPacketTime int64) {
	for _, s := range segments {
		if s.ContainerPath == "" {
			p.logger.Warnf("Audio segment %s was not converted, skipping", s.metadata.BaseFilename)
			continue
		}

		ts := firstPacketTimestamp(s.metadata)
		offsets = append(offsets, &FileOffset{Name: s.ContainerPath, Offset: ts})
		if firstPacketTime == 0 || ts < firstPacketTime {
			firstPacketTime = ts
		}
	}

	return offsets, firstPacketTime
}

func (p *AudioMixerV2) buildFilename(config *AudioMixerConfig, metadata *RecordingMetadata) string {
	var callType, callID string
	var startMs int64
	if len(metadata.Tracks) > 0 {
		tr := metadata.Tracks[0]
		callType, callID, startMs = tr.CallType, tr.CallID, tr.CallStartTime.UTC().UnixMilli()
	} else if metadata.CallMetadata != nil {
		cm := metadata.CallMetadata
		callType, callID, startMs = cm.CallType, cm.CallID, cm.RecordingStartTime.UTC().UnixMilli()
	}
	return filepath.Join(config.OutputDir, fmt.Sprintf("rec_%s_%s_%s_%d.%s", callType, callID, trackKindAudio, startMs, config.Format))
}
