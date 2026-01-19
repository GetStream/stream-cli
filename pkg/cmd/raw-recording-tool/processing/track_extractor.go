package processing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GetStream/getstream-go/v3"
	"github.com/pion/webrtc/v4"
)

// Generic track extraction function that works for both audio and video
func ExtractTracks(workingDir, outputDir, userID, sessionID, trackID string, metadata *RecordingMetadata, trackType, mediaFilter string, fillGaps, fixDtx bool, logger *getstream.DefaultLogger) error {
	// Filter tracks to specified type only and apply hierarchical filtering
	filteredTracks := FilterTracks(metadata.Tracks, userID, sessionID, trackID, trackType, mediaFilter)
	if len(filteredTracks) == 0 {
		logger.Warn("No %s tracks found matching the filter criteria", trackType)
		return nil
	}

	logger.Info("Found %d %s tracks to extract", len(filteredTracks), trackType)

	// Extract and convert each track
	for i, track := range filteredTracks {
		logger.Info("Processing %s track %d/%d: %s", trackType, i+1, len(filteredTracks), track.TrackID)

		err := extractSingleTrackWithOptions(workingDir, track, outputDir, trackType, fillGaps, fixDtx, logger)
		if err != nil {
			logger.Error("Failed to extract %s track %s: %v", trackType, track.TrackID, err)
			continue
		}
	}

	return nil
}

func extractSingleTrackWithOptions(inputPath string, track *TrackInfo, outputDir string, trackType string, fillGaps, fixDtx bool, logger *getstream.DefaultLogger) error {
	accept := func(path string, info os.FileInfo) (*SegmentInfo, bool) {
		for _, s := range track.Segments {
			if strings.Contains(info.Name(), s.metadata.BaseFilename) {
				if track.Codec == webrtc.MimeTypeH264 {
					s.ContainerExt = Mp4
				} else {
					s.ContainerExt = Webm
				}
				s.RtpDumpPath = path
				s.SdpPath = strings.Replace(path, SuffixRtpDump, SuffixSdp, -1)
				s.ContainerPath = strings.Replace(path, SuffixRtpDump, "."+s.ContainerExt, -1)
				return s, true
			}
		}
		return nil, false
	}

	// Convert using the WebM converter
	err := ConvertDirectory(inputPath, accept, fixDtx, logger)
	if err != nil {
		return fmt.Errorf("failed to convert %s track: %w", trackType, err)
	}

	// Create segments with timing info and fill gaps
	finalFile, err := processSegmentsWithGapFilling(track, trackType, outputDir, fillGaps, logger)
	if err != nil {
		return fmt.Errorf("failed to process segments with gap filling: %w", err)
	}

	track.ConcatenatedContainerPath = finalFile
	logger.Info("Successfully extracted %s track to: %s", trackType, finalFile)
	return nil
}

// processSegmentsWithGapFilling processes webm segments, fills gaps if requested, and concatenates into final file
func processSegmentsWithGapFilling(track *TrackInfo, trackType string, outputDir string, fillGaps bool, logger *getstream.DefaultLogger) (string, error) {
	// Build list of files to concatenate (with optional gap fillers)
	var filesToConcat []string
	for i, segment := range track.Segments {
		// Add the segment file
		filesToConcat = append(filesToConcat, segment.ContainerPath)

		// Add gap filler if requested and there's a gap before the next segment
		if fillGaps && i < track.SegmentCount-1 {
			nextSegment := track.Segments[i+1]
			gapDuration := nextSegment.FFMpegOffset + firstPacketNtpTimestamp(nextSegment.metadata) - lastPacketNtpTimestamp(segment.metadata)

			if gapDuration > 0 { // There's a gap
				gapSeconds := float64(gapDuration) / 1000.0
				logger.Info("Detected %dms gap between segments, generating %s filler", gapDuration, trackType)

				// Create gap filler file
				gapFilePath := filepath.Join(outputDir, fmt.Sprintf("gap_%s_%d.%s", trackType, i, segment.ContainerExt))

				if trackType == "audio" {
					err := generateSilence(gapFilePath, gapSeconds, logger)
					if err != nil {
						logger.Warn("Failed to generate silence, skipping gap: %v", err)
						continue
					}
				} else if trackType == "video" {
					// Use 720p quality as defaults
					err := generateBlackVideo(gapFilePath, track.Codec, gapSeconds, 1280, 720, 30, logger)
					if err != nil {
						logger.Warn("Failed to generate black video, skipping gap: %v", err)
						continue
					}
				}

				defer os.Remove(gapFilePath)

				filesToConcat = append(filesToConcat, gapFilePath)
			}
		}
	}

	// Create final output file
	finalName := fmt.Sprintf("%s_%s_%s_%s.%s", trackType, track.UserID, track.SessionID, track.TrackID, track.Segments[0].ContainerExt)
	finalPath := filepath.Join(outputDir, finalName)

	// Concatenate all segments (with gap fillers if any)
	err := concatFile(finalPath, filesToConcat, logger)
	if err != nil {
		return "", fmt.Errorf("failed to concatenate segments: %w", err)
	}

	logger.Info("Successfully concatenated %d segments into %s (gap filled %t)", track.SegmentCount, finalPath, fillGaps)
	return finalPath, nil
}
