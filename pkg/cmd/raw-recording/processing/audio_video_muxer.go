package processing

import (
	"errors"
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

func (p *AudioVideoMuxer) MuxAudioVideoTracks(config *AudioVideoMuxerConfig, metadata *RecordingMetadata, extracted []*TrackFileInfo) ([]*TrackFileInfo, error) {
	if config.WithExtract || extracted == nil {
		cfg := &TrackExtractorConfig{
			WorkDir:        config.WorkDir,
			OutputDir:      config.OutputDir,
			UserID:         config.UserID,
			SessionID:      config.SessionID,
			TrackID:        config.TrackID,
			TrackKind:      "",
			MediaType:      config.MediaType,
			AggregatorType: AggregatorTypeSession,
			FillGap:        true,
			FillDtx:        true,

			Cleanup: config.WithCleanup,
		}

		extractor := NewTrackExtractor(p.logger)

		// Extract tracks with gap filling enabled
		p.logger.Infof("Extracting tracks with gap filling...")
		tracks, err := extractor.ExtractTracks(cfg, metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to extract audio tracks: %w", err)
		}
		extracted = tracks
	}

	var infos []*TrackFileInfo // Group files by media type for proper pairing
	pairedTracks := p.groupFilesByMediaType(config, extracted)

	for audioTrack, videoTrack := range pairedTracks {
		info, err := p.muxTrackPairs(audioTrack, videoTrack, config)
		if err != nil {
			p.logger.Errorf("Failed to mux user tracks: %v", err)
			continue
		}
		infos = append(infos, info)
	}

	return infos, nil
}

// calculateSyncOffsetFromFiles calculates sync offset between audio and video files using metadata
func calculateSyncOffsetFromFiles(firstFile, secondFile *TrackFileInfo) (int64, error) {
	// Calculate offset: positive means video starts before audio
	var firstTrack, secondTrack *TrackInfo
	if len(firstFile.AudioTrack) > 0 {
		firstTrack = firstFile.AudioTrack[0]
	} else if len(firstFile.VideoTrack) > 0 {
		firstTrack = firstFile.VideoTrack[0]
	} else {
		return 0, errors.New("no track found in first file")
	}

	if len(secondFile.AudioTrack) > 0 {
		secondTrack = secondFile.AudioTrack[0]
	} else if len(secondFile.VideoTrack) > 0 {
		secondTrack = secondFile.VideoTrack[0]
	} else {
		return 0, errors.New("no track found in second file")
	}

	if len(firstTrack.Segments) == 0 {
		return 0, errors.New("no segment found in first track")
	}

	if len(secondTrack.Segments) == 0 {
		return 0, errors.New("no segment found in second track")
	}

	firstTs := firstPacketTimestamp(firstTrack.Segments[0].metadata)
	secondTs := firstPacketTimestamp(secondTrack.Segments[0].metadata)
	offset := firstTs - secondTs

	return offset, nil
}

// groupFilesByMediaType groups audio and video files by media type (user vs display)
func (p *AudioVideoMuxer) groupFilesByMediaType(config *AudioVideoMuxerConfig, extracted []*TrackFileInfo) map[*TrackFileInfo]*TrackFileInfo {
	pairedTracks := make(map[*TrackFileInfo]*TrackFileInfo)

	matches := func(audio *TrackFileInfo, video *TrackFileInfo) bool {
		return audio.AudioTrack[0].UserID == video.VideoTrack[0].UserID &&
			audio.AudioTrack[0].SessionID == video.VideoTrack[0].SessionID &&
			audio.AudioTrack[0].IsScreenshare == video.VideoTrack[0].IsScreenshare
	}

	for _, at := range extracted {
		audioSize := len(at.AudioTrack)
		if audioSize > 0 && audioSize == len(FilterTracks(at.AudioTrack, config.UserID, config.SessionID, config.TrackID, "", config.MediaType)) {
			for _, vt := range extracted {
				videoSize := len(vt.VideoTrack)
				if videoSize > 0 && videoSize == len(FilterTracks(vt.VideoTrack, config.UserID, config.SessionID, config.TrackID, "", config.MediaType)) && matches(at, vt) {
					pairedTracks[at] = vt
					break
				}
			}
		}
	}

	return pairedTracks
}

// muxTrackPairs muxes audio/video pairs of the same media type
func (p *AudioVideoMuxer) muxTrackPairs(audio, video *TrackFileInfo, config *AudioVideoMuxerConfig) (*TrackFileInfo, error) {
	// Calculate sync offset using segment timing information
	offset, err := calculateSyncOffsetFromFiles(audio, video)
	if err != nil {
		p.logger.Warnf("Failed to calculate sync offset, using 0: %v", err)
		offset = 0
	}

	// Generate output filename with media type indicator
	outputFile := p.buildFilename(config.OutputDir, video)

	audioFile := audio.Name
	videoFile := video.Name

	// Mux the audio and video files
	p.logger.Debugf("Muxing %s + %s → %s (offset: %dms)",
		filepath.Base(audioFile), filepath.Base(videoFile), filepath.Base(outputFile), offset)

	err = runFFmpegCommand(generateMuxFilesArguments(outputFile, audioFile, videoFile, float64(offset)), p.logger)
	if err != nil {
		p.logger.Errorf("Failed to mux %s + %s: %v", audioFile, videoFile, err)
		return nil, err
	}

	p.logger.Infof("Successfully created muxed file: %s", outputFile)

	// Clean up individual track files to avoid clutter
	if config.WithCleanup {
		defer func() {
			for _, file := range []string{audioFile, videoFile} {
				p.logger.Infof("Cleaning up temporary file: %s", file)
				if err := os.Remove(file); err != nil {
					p.logger.Warnf("Failed to clean up temporary file %s: %v", file, err)
				}
			}
		}()
	}

	return &TrackFileInfo{
		Name:              outputFile,
		StartAt:           p.getTime(audio.StartAt, video.StartAt, true),
		EndAt:             p.getTime(audio.EndAt, video.EndAt, false),
		MaxFrameDimension: video.MaxFrameDimension,
		AudioTrack:        audio.AudioTrack,
		VideoTrack:        video.VideoTrack,
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
func (p *AudioVideoMuxer) buildFilename(outputDir string, file *TrackFileInfo) string {
	callType := "callType"
	callID := "callID"
	userID := "userId"
	sessionID := "sessionId"
	startTime := int64(0)
	ext := mkvExtension
	media := "audio_video"

	if len(file.VideoTrack) > 0 {
		track := file.VideoTrack[0]
		callType = track.CallType
		callID = track.CallID
		userID = track.UserID
		sessionID = track.SessionID
		if track.IsScreenshare {
			media = "shared_" + media
		}
		startTime = track.CallStartTime.UnixMilli()
	}

	return filepath.Join(outputDir, fmt.Sprintf("individual_%s_%s_%s_%s_%s_%d.%s", callType, callID, userID, sessionID, media, startTime, ext))
}
