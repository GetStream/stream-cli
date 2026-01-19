package processing

import (
	"fmt"
	"path/filepath"

	"github.com/GetStream/getstream-go/v3"
)

type AudioMixerConfig struct {
	WorkDir         string
	OutputDir       string
	WithScreenshare bool
	WithExtract     bool
	WithCleanup     bool
}

type AudioMixer struct {
	logger *getstream.DefaultLogger
}

func NewAudioMixer(logger *getstream.DefaultLogger) *AudioMixer {
	return &AudioMixer{logger: logger}
}

// MixAllAudioTracks orchestrates the entire audio mixing workflow using existing extraction logic
func (p *AudioMixer) MixAllAudioTracks(config *AudioMixerConfig, metadata *RecordingMetadata, logger *getstream.DefaultLogger) error {
	// Step 1: Extract all matching audio tracks using existing ExtractTracks function
	logger.Info("Step 1/2: Extracting all matching audio tracks...")

	if config.WithExtract {
		mediaFilter := "user"
		if config.WithScreenshare {
			mediaFilter = "both"
		}

		if err := ExtractTracks(config.WorkDir, config.OutputDir, "", "", "", metadata, "audio", mediaFilter, true, true, logger); err != nil {
			return fmt.Errorf("failed to extract audio tracks: %w", err)
		}
	}

	fileOffsetMap := p.offset(metadata, config.WithScreenshare, logger)
	if len(fileOffsetMap) == 0 {
		return fmt.Errorf("no audio files were extracted - check your filter criteria")
	}

	logger.Info("Found %d extracted audio files to mix", len(fileOffsetMap))

	// Step 3: Mix all discovered audio files using existing webm.mixAudioFiles
	outputFile := filepath.Join(config.OutputDir, "mixed_audio.webm")

	err := mixAudioFiles(outputFile, fileOffsetMap, logger)
	if err != nil {
		return fmt.Errorf("failed to mix audio files: %w", err)
	}

	logger.Info("Successfully created mixed audio file: %s", outputFile)

	//// Clean up individual audio files (optional)
	//for _, audioFile := range audioFiles {
	//	if err := os.Remove(audioFile.FilePath); err != nil {
	//		logger.Warn("Failed to clean up temporary file %s: %v", audioFile.FilePath, err)
	//	}
	//}

	return nil
}

func (p *AudioMixer) offset(metadata *RecordingMetadata, withScreenshare bool, logger *getstream.DefaultLogger) []*FileOffset {
	var offsets []*FileOffset
	var firstTrack *TrackInfo
	for _, t := range metadata.Tracks {
		if t.TrackType == "audio" && (!t.IsScreenshare || withScreenshare) {
			if firstTrack == nil {
				firstTrack = t
				offsets = append(offsets, &FileOffset{
					Name:   t.ConcatenatedContainerPath,
					Offset: 0, // Will be sorted later and rearranged
				})
			} else {
				offset, err := calculateSyncOffsetFromFiles(t, firstTrack, logger)
				if err != nil {
					logger.Warn("Failed to calculate sync offset for audio tracks: %v", err)
					continue
				}

				offsets = append(offsets, &FileOffset{
					Name:   t.ConcatenatedContainerPath,
					Offset: offset,
				})
			}
		}
	}

	return offsets
}
