package processing

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"
)

const (
	FormatMp3     = "mp3"
	FormatWeba    = "weba"
	FormatWebm    = "webm"
	FormatMka     = "mka"
	FormatMkv     = "mkv"
	DefaultFormat = FormatMkv
)

var supportedFormats = [5]string{FormatMp3, FormatWeba, FormatWebm, FormatMka, FormatMkv}

type AudioMixerConfig struct {
	WorkDir         string
	OutputDir       string
	Format          string
	WithScreenshare bool

	WithExtract bool
	WithCleanup bool

	// RecordingStart and RecordingEnd define the recording window. When both
	// are set and RecordingEnd is after RecordingStart, the mixer anchors a
	// non-empty output to the window without generating any silence files: each track is
	// delayed (adelay) by its offset from RecordingStart to produce leading
	// silence, and the mix is padded (apad) to the full window duration to
	// produce trailing silence. When no participant audio exists, no output is
	// generated. Leaving them zero preserves legacy behaviour (output spans only
	// the captured tracks).
	RecordingStart time.Time
	RecordingEnd   time.Time
}

type AudioMixer struct {
	logger *ProcessingLogger
}

func NewAudioMixer(logger *ProcessingLogger) *AudioMixer {
	return &AudioMixer{logger: logger}
}

// MixAllAudioTracks orchestrates the entire audio mixing workflow using existing extraction logic
func (p *AudioMixer) MixAllAudioTracks(config *AudioMixerConfig, metadata *RecordingMetadata, extracted []*TrackFileInfo) (*string, error) {
	p.overrideConfig(config)

	// Step 1: Extract all matching audio tracks using existing ExtractTracks function
	p.logger.Infof("Extracting all matching audio tracks...")

	if config.WithExtract || extracted == nil {
		mediaType := ""
		if !config.WithScreenshare {
			mediaType = "user"
		}

		cfg := &TrackExtractorConfig{
			WorkDir:        config.WorkDir,
			OutputDir:      config.OutputDir,
			UserID:         "",
			SessionID:      "",
			TrackID:        "",
			TrackKind:      trackKindAudio,
			MediaType:      mediaType,
			AggregatorType: AggregatorTypeTrack,
			FillDtx:        true,
			FillGap:        true,

			Cleanup: config.WithCleanup,
		}

		extractor := NewTrackExtractor(p.logger)
		files, err := extractor.ExtractTracks(cfg, metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to extract audio tracks: %w", err)
		}
		extracted = files
	}

	hasWindow := config.RecordingEnd.After(config.RecordingStart)

	fileOffsets, firstPacketTime := p.offset(extracted, config.WithScreenshare)
	if len(fileOffsets) == 0 {
		p.logger.Warnf("No audio tracks found; skipping composite output")
		return nil, nil
	}

	p.logger.Infof("Found %d extracted audio files to mix", len(fileOffsets))

	//// Clean up individual audio files (optional)
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

	// Step 3: Mix all discovered audio files using existing webm.mixAudioFiles
	outputFile := p.buildFilename(config, metadata)

	err := runFFmpegCommand(generateMixAudioFilesArguments(outputFile, config.Format, fileOffsets, leadingGapMs, recordingDurationSeconds), p.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to mix audio files: %w", err)
	}

	p.logger.Infof("Successfully created mixed audio file: %s", outputFile)

	return &outputFile, nil
}

func (p *AudioMixer) overrideConfig(config *AudioMixerConfig) {
	if !slices.Contains(supportedFormats[:], config.Format) {
		p.logger.Warnf("Audio format %s not supported, fallback to default %s", config.Format, DefaultFormat)
		config.Format = DefaultFormat
	}
}

func (p *AudioMixer) offset(extracted []*TrackFileInfo, withScreenshare bool) (offsets []*FileOffset, firstPacketTime int64) {
	var firstTrack *TrackFileInfo
	for _, t := range extracted {
		if len(t.AudioTrack) == 0 || (t.AudioTrack[0].IsScreenshare && !withScreenshare) {
			continue
		}

		if firstTrack == nil {
			firstTrack = t
			offsets = append(offsets, &FileOffset{
				Name:   t.Name,
				Offset: 0, // Will be sorted later and rearranged
			})
		} else {
			offset, err := calculateSyncOffsetFromFiles(t, firstTrack)
			if err != nil {
				p.logger.Warnf("Failed to calculate sync offset for audio tracks: %v", err)
				continue
			}

			offsets = append(offsets, &FileOffset{
				Name:   t.Name,
				Offset: offset,
			})
		}

		firstPktTs := firstPacketTimestamp(t.AudioTrack[0].Segments[0].metadata)
		if firstPacketTime == 0 || firstPktTs < firstPacketTime {
			firstPacketTime = firstPktTs
		}
	}

	return offsets, firstPacketTime
}

func (p *AudioMixer) buildFilename(config *AudioMixerConfig, metadata *RecordingMetadata) string {
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
