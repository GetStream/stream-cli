package processing

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const blackVideoFps = 5

type TrackExtractorConfig struct {
	WorkDir   string
	OutputDir string
	UserID    string
	SessionID string
	TrackID   string
	TrackKind string
	MediaType string
	FillGap   bool
	FillDtx   bool

	Cleanup bool
}

type TrackExtractor struct {
	logger *ProcessingLogger
}

func NewTrackExtractor(logger *ProcessingLogger) *TrackExtractor {
	return &TrackExtractor{logger: logger}
}

// Generic track extraction function that works for both audio and video
func (p *TrackExtractor) ExtractTracks(config *TrackExtractorConfig, metadata *RecordingMetadata) ([]*TrackFileInfo, error) {
	// Filter tracks to specified type only and apply hierarchical filtering
	filteredTracks := FilterTracks(metadata.Tracks, config.UserID, config.SessionID, config.TrackID, config.TrackKind, config.MediaType)
	if len(filteredTracks) == 0 {
		p.logger.Warnf("No %s tracks found matching the filter criteria", config.TrackKind)
		return nil, nil
	}

	p.logger.Infof("Found %d %s tracks to extract", len(filteredTracks), config.TrackKind)

	// Extract and convert each track
	var infos []*TrackFileInfo
	for i, track := range filteredTracks {
		p.logger.Debugf("Processing %s track %d/%d: %s", track.TrackKind, i+1, len(filteredTracks), track.TrackID)

		info, err := p.extractSingleTrackWithOptions(config, track)
		if err != nil {
			p.logger.Errorf("Failed to extract %s track %s: %v", track.TrackKind, track.TrackID, err)
			continue
		}
		if info != nil {
			infos = append(infos, info)
		}
	}

	return infos, nil
}

func (p *TrackExtractor) extractSingleTrackWithOptions(config *TrackExtractorConfig, track *TrackInfo) (*TrackFileInfo, error) {
	accept := func(path string, info os.FileInfo) (*SegmentInfo, bool) {
		for _, s := range track.Segments {
			if strings.Contains(info.Name(), s.metadata.BaseFilename) {
				extension, suffix := outputFormatForMimeType(track.Codec)
				abs, _ := filepath.Abs(path)

				s.RtpDumpPath = abs
				s.SdpPath = strings.Replace(abs, suffixRtpDump, suffixSdp, -1)
				s.ContainerExt = extension
				s.ContainerPath = strings.Replace(abs, suffixRtpDump, suffix, -1)
				return s, true
			}
		}
		return nil, false
	}

	// Convert using the WebM converter
	err := ConvertDirectory(config.WorkDir, accept, config.FillDtx, p.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to convert %s track: %w", track.TrackKind, err)
	}

	// Create segments with timing info and fill gaps
	finalFileInfo, err := p.processSegmentsWithGapFilling(config, track)
	if err != nil {
		return nil, fmt.Errorf("failed to process segments with gap filling: %w", err)
	}

	track.ConcatenatedTrackFileInfo = finalFileInfo
	p.logger.Infof("Successfully extracted %s track to: %s", track.TrackKind, finalFileInfo.Name)
	return finalFileInfo, nil
}

// processSegmentsWithGapFilling processes webm segments, fills gaps if requested, and concatenates into final file
func (p *TrackExtractor) processSegmentsWithGapFilling(config *TrackExtractorConfig, track *TrackInfo) (*TrackFileInfo, error) {
	// Build list of files to concatenate (with optional gap fillers)
	var cleanupFiles []string
	concatFile, err := os.Create(p.buildConcatFilename(config.OutputDir, track))
	if err != nil {
		return nil, err
	}
	cleanupFiles = append(cleanupFiles, concatFile.Name())

	// If enabled, cleanUp all working files (segment mkv, silence or black frame files and concat.txt)
	if config.Cleanup {
		defer func(files *[]string) {
			for _, file := range *files {
				p.logger.Infof("Cleaning up temporary file: %s", file)
				if err := os.Remove(file); err != nil {
					p.logger.Warnf("Failed to clean up temporary file %s: %v", file, err)
				}
			}
		}(&cleanupFiles)
	}
	defer concatFile.Close()

	for i, segment := range track.Segments {
		if _, e := concatFile.WriteString(fmt.Sprintf("file '%s'\n", segment.ContainerPath)); e != nil {
			return nil, e
		}
		cleanupFiles = append(cleanupFiles, segment.ContainerPath)

		// Add gap filler if requested and there's a gap before the next segment
		if config.FillGap && i < track.SegmentCount-1 {
			nextSegment := track.Segments[i+1]
			offset := int64(0)
			if nextSegment.metadata.FirstKeyFrameOffsetMs != nil {
				offset = *nextSegment.metadata.FirstKeyFrameOffsetMs
			}
			gapDuration := offset + firstPacketNtpTimestamp(nextSegment.metadata) - lastPacketNtpTimestamp(segment.metadata)

			if gapDuration > 0 { // There's a gap
				gapSeconds := float64(gapDuration) / 1000.0
				p.logger.Infof("Detected %dms gap between segments, generating %s filler", gapDuration, track.TrackKind)

				// Create gap filler file
				gapFilePath := p.buildGapFilename(config.OutputDir, track, i)

				var args []string
				if track.TrackKind == trackKindVideo {
					args = generateBlackVideoArguments(gapFilePath, track.Codec, gapSeconds, 1280, 720, blackVideoFps)
				} else {
					args = generateSilenceArguments(gapFilePath, gapSeconds)
				}

				if e := runFFmpegCommand(args, p.logger); e != nil {
					p.logger.Warnf("Failed to generate %s gap, skipping: %v", track.TrackKind, e)
					continue
				}
				cleanupFiles = append(cleanupFiles, gapFilePath)

				absPath, err := filepath.Abs(gapFilePath)
				if err != nil {
					return nil, err
				}

				if _, e := concatFile.WriteString(fmt.Sprintf("file '%s'\n", absPath)); e != nil {
					return nil, e
				}
			}
		}
	}

	// Create final output file
	finalPath := p.buildFilename(config.OutputDir, track)

	// Concatenate all segments (with gap fillers if any)
	args, err := generateConcatFileArguments(finalPath, concatFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to generate ffmpeg arguments: %w", err)
	}

	err = runFFmpegCommand(args, p.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to concatenate segments: %w", err)
	}

	p.logger.Debugf("Successfully concatenated %d segments into %s (gap filled %t)", track.SegmentCount, finalPath, config.FillGap)

	var ts, te int64
	if len(track.Segments) > 0 {
		ts = track.Segments[0].metadata.FirstRtpUnixTimestamp
		te = track.Segments[len(track.Segments)-1].metadata.LastRtpUnixTimestamp
	}

	var audioTrack, videoTrack *TrackInfo
	switch track.TrackKind {
	case trackKindAudio:
		audioTrack = track
	case trackKindVideo:
		videoTrack = track
	}
	return &TrackFileInfo{
		Name:              finalPath,
		StartAt:           time.UnixMilli(ts),
		EndAt:             time.UnixMilli(te),
		MaxFrameDimension: p.getMaxFrameDimension(track),
		AudioTrack:        audioTrack,
		VideoTrack:        videoTrack,
	}, nil
}

func (p *TrackExtractor) getMaxFrameDimension(track *TrackInfo) SegmentFrameDimension {
	frameDimension := SegmentFrameDimension{}
	if track.TrackKind == trackKindVideo {
		for _, segment := range track.Segments {
			if segment.metadata.MaxFrameDimension != nil {
				frameDimension = getMaxFrameDimension(*segment.metadata.MaxFrameDimension, frameDimension)
			}
		}
	}
	return frameDimension
}

// buildDefaultFilename creates output filename that indicates media type
func (p *TrackExtractor) buildFilename(outputDir string, track *TrackInfo) string {
	media := track.TrackKind + "_only"
	if track.IsScreenshare {
		media = "shared_" + media
	}

	return filepath.Join(outputDir, fmt.Sprintf("individual_%s_%s_%s_%s_%s_%d.%s", track.CallType, track.CallID, track.UserID, track.SessionID, media, track.CallStartTime.UnixMilli(), track.Segments[0].ContainerExt))
}

func (p *TrackExtractor) buildGapFilename(outputDir string, track *TrackInfo, i int) string {
	return filepath.Join(outputDir, fmt.Sprintf("gap_%s_%s_%s_%d_%d.%s", track.UserID, track.SessionID, track.TrackKind, track.CallStartTime.UnixMilli(), i, track.Segments[i].ContainerExt))
}

func (p *TrackExtractor) buildConcatFilename(outputDir string, track *TrackInfo) string {
	return filepath.Join(outputDir, fmt.Sprintf("concat_%s_%s_%s_%d.txt", track.UserID, track.SessionID, track.TrackKind, track.CallStartTime.UnixMilli()))
}
