package processing

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/GetStream/getstream-go/v3"
)

type AudioVideoMuxerConfig struct {
	WorkDir   string
	OutputDir string
	UserID    string
	SessionID string
	TrackID   string
	Media     string

	WithExtract bool
	WithCleanup bool
}

type AudioVideoMuxer struct {
	logger *getstream.DefaultLogger
}

func NewAudioVideoMuxer(logger *getstream.DefaultLogger) *AudioVideoMuxer {
	return &AudioVideoMuxer{logger: logger}
}

func (p *AudioVideoMuxer) MuxAudioVideoTracks(config *AudioVideoMuxerConfig, metadata *RecordingMetadata, logger *getstream.DefaultLogger) error {
	if config.WithExtract {
		// Extract audio tracks with gap filling enabled
		logger.Info("Extracting audio tracks with gap filling...")
		err := ExtractTracks(config.WorkDir, config.OutputDir, config.UserID, config.SessionID, config.TrackID, metadata, "audio", config.Media, true, true, logger)
		if err != nil {
			return fmt.Errorf("failed to extract audio tracks: %w", err)
		}

		// Extract video tracks with gap filling enabled
		logger.Info("Extracting video tracks with gap filling...")
		err = ExtractTracks(config.WorkDir, config.OutputDir, config.UserID, config.SessionID, config.TrackID, metadata, "video", config.Media, true, true, logger)
		if err != nil {
			return fmt.Errorf("failed to extract video tracks: %w", err)
		}
	}

	// Group files by media type for proper pairing
	pairedTracks := p.groupFilesByMediaType(metadata)
	for audioTrack, videoTrack := range pairedTracks {
		//logger.Info("Muxing %d user audio/video pairs", len(userAudio))
		err := p.muxTrackPairs(audioTrack, videoTrack, config.OutputDir, logger)
		if err != nil {
			logger.Error("Failed to mux user tracks: %v", err)
		}
	}

	return nil
}

// calculateSyncOffsetFromFiles calculates sync offset between audio and video files using metadata
func calculateSyncOffsetFromFiles(audioTrack, videoTrack *TrackInfo, logger *getstream.DefaultLogger) (int64, error) {
	// Calculate offset: positive means video starts before audio
	audioTs := audioTrack.Segments[0].FFMpegOffset + firstPacketNtpTimestamp(audioTrack.Segments[0].metadata)
	videoTs := videoTrack.Segments[0].FFMpegOffset + firstPacketNtpTimestamp(videoTrack.Segments[0].metadata)
	offset := audioTs - videoTs

	logger.Info(fmt.Sprintf("Calculated sync offset: audio_start=%v, audio_ts=%v, video_start=%v, video_ts=%v, offset=%d",
		audioTrack.Segments[0].metadata.FirstRtpUnixTimestamp, audioTs, videoTrack.Segments[0].metadata.FirstRtpUnixTimestamp, videoTs, offset))

	return offset, nil
}

// groupFilesByMediaType groups audio and video files by media type (user vs display)
func (p *AudioVideoMuxer) groupFilesByMediaType(metadata *RecordingMetadata) map[*TrackInfo]*TrackInfo {
	pairedTracks := make(map[*TrackInfo]*TrackInfo)

	matches := func(audio *TrackInfo, video *TrackInfo) bool {
		return audio.UserID == video.UserID &&
			audio.SessionID == video.SessionID &&
			audio.IsScreenshare == video.IsScreenshare
	}

	for _, at := range metadata.Tracks {
		if at.TrackType == "audio" {
			for _, vt := range metadata.Tracks {
				if vt.TrackType == "video" && matches(at, vt) {
					pairedTracks[at] = vt
					break
				}
			}
		}
	}

	return pairedTracks
}

// muxTrackPairs muxes audio/video pairs of the same media type
func (p *AudioVideoMuxer) muxTrackPairs(audio, video *TrackInfo, outputDir string, logger *getstream.DefaultLogger) error {
	// Calculate sync offset using segment timing information
	offset, err := calculateSyncOffsetFromFiles(audio, video, logger)
	if err != nil {
		logger.Warn("Failed to calculate sync offset, using 0: %v", err)
		offset = 0
	}

	// Generate output filename with media type indicator
	outputFile := p.generateMediaAwareMuxedFilename(audio, video, outputDir)

	audioFile := audio.ConcatenatedContainerPath
	videoFile := video.ConcatenatedContainerPath

	// Mux the audio and video files
	logger.Info("Muxing %s + %s → %s (offset: %dms)",
		filepath.Base(audioFile), filepath.Base(videoFile), filepath.Base(outputFile), offset)

	err = muxFiles(outputFile, audioFile, videoFile, float64(offset), logger)
	if err != nil {
		logger.Error("Failed to mux %s + %s: %v", audioFile, videoFile, err)
		return err
	}

	logger.Info("Successfully created muxed file: %s", outputFile)

	// Clean up individual track files to avoid clutter
	//os.Remove(audioFile)
	//os.Remove(videoFile)
	//}
	//
	//if len(audioFiles) != len(videoFiles) {
	//	logger.Warn("Mismatched %s track counts: %d audio, %d video", mediaTypeName, len(audioFiles), len(videoFiles))
	//}

	return nil
}

// generateMediaAwareMuxedFilename creates output filename that indicates media type
func (p *AudioVideoMuxer) generateMediaAwareMuxedFilename(audioFile, videoFile *TrackInfo, outputDir string) string {
	audioBase := filepath.Base(audioFile.Segments[0].ContainerPath)
	audioBase = strings.TrimSuffix(audioBase, "."+audioFile.Segments[0].ContainerExt)

	// Replace "audio_" with "muxed_{mediaType}_" to create output name
	var muxedName string
	if audioFile.IsScreenshare {
		muxedName = strings.Replace(audioBase, "audio_", "muxed_display_", 1) + "." + videoFile.Segments[0].ContainerExt
	} else {
		muxedName = strings.Replace(audioBase, "audio_", "muxed_", 1) + "." + videoFile.Segments[0].ContainerExt
	}

	return filepath.Join(outputDir, muxedName)
}
