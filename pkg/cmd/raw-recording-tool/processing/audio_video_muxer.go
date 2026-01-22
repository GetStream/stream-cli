package processing

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type AudioVideoMuxerConfig struct {
	WorkDir   string
	OutputDir string
	UserID    string
	SessionID string
	TrackID   string
	MediaType string

	WithExtract bool
	WithCleanup bool
}

type AudioVideoMuxer struct {
	logger *ProcessingLogger
}

func NewAudioVideoMuxer(logger *ProcessingLogger) *AudioVideoMuxer {
	return &AudioVideoMuxer{logger: logger}
}

func (p *AudioVideoMuxer) MuxAudioVideoTracks(config *AudioVideoMuxerConfig, metadata *RecordingMetadata) ([]*TrackFileInfo, error) {
	if config.WithExtract {
		cfg := &TrackExtractorConfig{
			WorkDir:   config.WorkDir,
			OutputDir: config.OutputDir,
			UserID:    config.UserID,
			SessionID: config.SessionID,
			TrackID:   config.TrackID,
			TrackKind: "",
			MediaType: config.MediaType,
			FillGap:   true,
			FillDtx:   true,

			Cleanup: config.WithCleanup,
		}

		extractor := NewTrackExtractor(p.logger)

		// Extract tracks with gap filling enabled
		p.logger.Info("Extracting tracks with gap filling...")
		_, err := extractor.ExtractTracks(cfg, metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to extract audio tracks: %w", err)
		}
	}

	var infos []*TrackFileInfo // Group files by media type for proper pairing
	pairedTracks := p.groupFilesByMediaType(config, metadata)

	for audioTrack, videoTrack := range pairedTracks {
		// logger.Infof("Muxing %d user audio/video pairs", len(userAudio))
		info, err := p.muxTrackPairs(audioTrack, videoTrack, config)
		if err != nil {
			p.logger.Error("Failed to mux user tracks: %v", err)
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// calculateSyncOffsetFromFiles calculates sync offset between audio and video files using metadata
func calculateSyncOffsetFromFiles(audioTrack, videoTrack *TrackInfo) (int64, error) {
	// Calculate offset: positive means video starts before audio
	audioOffset, videoOffset := int64(0), int64(0)
	if audioTrack.Segments[0].metadata.FirstKeyFrameOffsetMs != nil {
		audioOffset = *audioTrack.Segments[0].metadata.FirstKeyFrameOffsetMs
	}
	if videoTrack.Segments[0].metadata.FirstKeyFrameOffsetMs != nil {
		videoOffset = *videoTrack.Segments[0].metadata.FirstKeyFrameOffsetMs
	}

	audioTs := audioOffset + firstPacketNtpTimestamp(audioTrack.Segments[0].metadata)
	videoTs := videoOffset + firstPacketNtpTimestamp(videoTrack.Segments[0].metadata)
	offset := audioTs - videoTs

	return offset, nil
}

// groupFilesByMediaType groups audio and video files by media type (user vs display)
func (p *AudioVideoMuxer) groupFilesByMediaType(config *AudioVideoMuxerConfig, metadata *RecordingMetadata) map[*TrackInfo]*TrackInfo {
	pairedTracks := make(map[*TrackInfo]*TrackInfo)

	matches := func(audio *TrackInfo, video *TrackInfo) bool {
		return audio.UserID == video.UserID &&
			audio.SessionID == video.SessionID &&
			audio.IsScreenshare == video.IsScreenshare
	}

	filteredTracks := FilterTracks(metadata.Tracks, config.UserID, config.SessionID, config.TrackID, "", config.MediaType)
	for _, at := range filteredTracks {
		if at.TrackKind == trackKindAudio {
			for _, vt := range filteredTracks {
				if vt.TrackKind == trackKindVideo && matches(at, vt) {
					pairedTracks[at] = vt
					break
				}
			}
		}
	}

	return pairedTracks
}

// muxTrackPairs muxes audio/video pairs of the same media type
func (p *AudioVideoMuxer) muxTrackPairs(audio, video *TrackInfo, config *AudioVideoMuxerConfig) (*TrackFileInfo, error) {
	// Calculate sync offset using segment timing information
	offset, err := calculateSyncOffsetFromFiles(audio, video)
	if err != nil {
		p.logger.Warn("Failed to calculate sync offset, using 0: %v", err)
		offset = 0
	}

	// Generate output filename with media type indicator
	outputFile := p.buildFilename(config.OutputDir, video)

	audioFile := audio.ConcatenatedTrackFileInfo.Name
	videoFile := video.ConcatenatedTrackFileInfo.Name

	// Mux the audio and video files
	p.logger.Debug("Muxing %s + %s → %s (offset: %dms)",
		filepath.Base(audioFile), filepath.Base(videoFile), filepath.Base(outputFile), offset)

	err = runFFmpegCommand(generateMuxFilesArguments(outputFile, audioFile, videoFile, float64(offset)), p.logger)
	if err != nil {
		p.logger.Errorf("Failed to mux %s + %s: %v", audioFile, videoFile, err)
		return nil, err
	}

	p.logger.Info("Successfully created muxed file: %s", outputFile)

	// Clean up individual track files to avoid clutter
	if config.WithCleanup {
		defer func() {
			for _, file := range []string{audioFile, videoFile} {
				p.logger.Info("Cleaning up temporary file: %s", file)
				if err := os.Remove(file); err != nil {
					p.logger.Warn("Failed to clean up temporary file %s: %v", file, err)
				}
			}
		}()
	}

	return &TrackFileInfo{
		Name:              outputFile,
		StartAt:           p.getTime(audio.ConcatenatedTrackFileInfo.StartAt, video.ConcatenatedTrackFileInfo.StartAt, true),
		EndAt:             p.getTime(audio.ConcatenatedTrackFileInfo.EndAt, video.ConcatenatedTrackFileInfo.EndAt, false),
		MaxFrameDimension: video.ConcatenatedTrackFileInfo.MaxFrameDimension,
	}, nil
}

func (p *AudioVideoMuxer) getTime(d1, d2 time.Time, first bool) time.Time {
	if d1.Before(d2) {
		if first {
			return d1
		} else {
			return d2
		}
	} else {
		if first {
			return d2
		} else {
			return d1
		}
	}
}

// buildFilename creates output filename that indicates media type
func (p *AudioVideoMuxer) buildFilename(outputDir string, track *TrackInfo) string {
	media := "audio_video"
	if track.IsScreenshare {
		media = "shared_" + media
	}

	return filepath.Join(outputDir, fmt.Sprintf("individual_%s_%s_%s_%s_%s_%d.%s", track.CallType, track.CallID, track.UserID, track.SessionID, media, track.CallStartTime.UnixMilli(), track.Segments[0].ContainerExt))
}
