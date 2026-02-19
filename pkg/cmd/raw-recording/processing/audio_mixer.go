package processing

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
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
}

type AudioMixer struct {
	logger *ProcessingLogger
}

func NewAudioMixer(logger *ProcessingLogger) *AudioMixer {
	return &AudioMixer{logger: logger}
}

// MixAllAudioTracks orchestrates the entire audio mixing workflow using existing extraction logic
func (p *AudioMixer) MixAllAudioTracks(config *AudioMixerConfig, metadata *RecordingMetadata) (*string, error) {
	p.overrideConfig(config)

	// Step 1: Extract all matching audio tracks using existing ExtractTracks function
	p.logger.Info("Extracting all matching audio tracks...")

	if config.WithExtract {
		mediaType := ""
		if !config.WithScreenshare {
			mediaType = "user"
		}

		cfg := &TrackExtractorConfig{
			WorkDir:   config.WorkDir,
			OutputDir: config.OutputDir,
			UserID:    "",
			SessionID: "",
			TrackID:   "",
			TrackKind: trackKindAudio,
			MediaType: mediaType,
			FillDtx:   true,
			FillGap:   true,

			Cleanup: config.WithCleanup,
		}

		extractor := NewTrackExtractor(p.logger)
		if _, err := extractor.ExtractTracks(cfg, metadata); err != nil {
			return nil, fmt.Errorf("failed to extract audio tracks: %w", err)
		}
	}

	fileOffsets := p.offset(metadata, config.WithScreenshare)
	if len(fileOffsets) == 0 {
		p.logger.Warn("No audio tracks found")
		return nil, nil
	}

	p.logger.Info("Found %d extracted audio files to mix", len(fileOffsets))

	//// Clean up individual audio files (optional)
	if config.WithCleanup {
		defer func(offsets *[]*FileOffset) {
			for _, fileOffset := range *offsets {
				p.logger.Info("Cleaning up temporary file: %s", fileOffset.Name)
				if err := os.Remove(fileOffset.Name); err != nil {
					p.logger.Warn("Failed to clean up temporary file %s: %v", fileOffset.Name, err)
				}
			}
		}(&fileOffsets)
	}

	// Step 3: Mix all discovered audio files using existing webm.mixAudioFiles
	outputFile := p.buildFilename(config, metadata)

	err := runFFmpegCommand(generateMixAudioFilesArguments(outputFile, config.Format, fileOffsets), p.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to mix audio files: %w", err)
	}

	p.logger.Info("Successfully created mixed audio file: %s", outputFile)

	return &outputFile, nil
}

func (p *AudioMixer) overrideConfig(config *AudioMixerConfig) {
	if !slices.Contains(supportedFormats[:], config.Format) {
		p.logger.Warn("Audio format %s not supported, fallback to default %s", config.Format, DefaultFormat)
		config.Format = DefaultFormat
	}
}

func (p *AudioMixer) offset(metadata *RecordingMetadata, withScreenshare bool) []*FileOffset {
	var offsets []*FileOffset
	var firstTrack *TrackInfo
	for _, t := range metadata.Tracks {
		if t.TrackKind == trackKindAudio && (!t.IsScreenshare || withScreenshare) {
			if firstTrack == nil {
				firstTrack = t
				offsets = append(offsets, &FileOffset{
					Name:   t.ConcatenatedTrackFileInfo.Name,
					Offset: 0, // Will be sorted later and rearranged
				})
			} else {
				offset, err := calculateSyncOffsetFromFiles(t, firstTrack)
				if err != nil {
					p.logger.Warn("Failed to calculate sync offset for audio tracks: %v", err)
					continue
				}

				offsets = append(offsets, &FileOffset{
					Name:   t.ConcatenatedTrackFileInfo.Name,
					Offset: offset,
				})
			}
		}
	}

	return offsets
}

func (p *AudioMixer) buildFilename(config *AudioMixerConfig, metadata *RecordingMetadata) string {
	tr := metadata.Tracks[0]
	return filepath.Join(config.OutputDir, fmt.Sprintf("composite_%s_%s_%s_%d.%s", tr.CallType, tr.CallID, trackKindAudio, tr.CallStartTime.UTC().UnixMilli(), config.Format))
}
