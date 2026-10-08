package processing

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const blackVideoFps = 5

type TrackExtractorConfig struct {
	WorkDir        string
	OutputDir      string
	UserID         string
	SessionID      string
	TrackID        string
	TrackKind      string
	MediaType      string
	AggregatorType string
	FillGap        bool
	FillDtx        bool

	Cleanup bool
}

// aggregationGroup holds TrackFileInfos that share the same session and track kind.
type aggregationGroup struct {
	TrackID       string
	SessionID     string
	UserID        string
	TrackKind     string
	IsScreenshare bool
	Infos         []*TrackInfo
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

	groups := p.groupTracks(filteredTracks, config.AggregatorType)

	// Extract and convert each track
	var infos []*TrackFileInfo
	i := 0
	for _, track := range groups {
		i++
		p.logger.Debugf("Processing %s track %d/%d: %s", track.TrackKind, i, len(groups), track.TrackID)

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

func (p *TrackExtractor) extractSingleTrackWithOptions(config *TrackExtractorConfig, track *aggregationGroup) (*TrackFileInfo, error) {
	accept := func(path string, info os.FileInfo) (*SegmentInfo, bool) {
		for _, i := range track.Infos {
			for _, s := range i.Segments {
				if strings.Contains(info.Name(), s.metadata.BaseFilename) {
					extension, suffix := outputFormatForMimeType(i.Codec)
					abs, _ := filepath.Abs(path)

					s.RtpDumpPath = abs
					s.SdpPath = strings.ReplaceAll(abs, suffixRtpDump, suffixSdp)
					s.ContainerExt = extension
					s.ContainerPath = strings.ReplaceAll(abs, suffixRtpDump, suffix)
					return s, true
				}
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

	p.logger.Infof("Successfully extracted %s track to: %s", track.TrackKind, finalFileInfo.Name)
	return finalFileInfo, nil
}

func (p *TrackExtractor) groupTracks(extractedInfos []*TrackInfo, aggregatorType string) []*aggregationGroup {
	groupMap := make(map[string]*aggregationGroup)

	for _, info := range extractedInfos {
		sessionID, trackID := info.SessionID, info.TrackID
		switch aggregatorType {
		case AggregatorTypeTrack: // Nothing to do, this is default
		case AggregatorTypeSession:
			trackID = ""
		}

		key := fmt.Sprintf("%s|%s|%s|%s|%t", info.UserID, sessionID, trackID, info.TrackKind, info.IsScreenshare)
		g, ok := groupMap[key]
		if !ok {
			g = &aggregationGroup{
				TrackID:       trackID,
				SessionID:     sessionID,
				UserID:        info.UserID,
				TrackKind:     info.TrackKind,
				IsScreenshare: info.IsScreenshare,
			}
			groupMap[key] = g
		}
		g.Infos = append(g.Infos, info)
	}

	groups := make([]*aggregationGroup, 0, len(groupMap))
	for _, group := range groupMap {
		slices.SortFunc(group.Infos, func(a, b *TrackInfo) int {
			return a.TrackStartTime.Compare(b.TrackStartTime)
		})
		groups = append(groups, group)
	}

	return groups
}

// processSegmentsWithGapFilling processes webm segments, fills gaps if requested, and concatenates into final file
func (p *TrackExtractor) processSegmentsWithGapFilling(config *TrackExtractorConfig, track *aggregationGroup) (*TrackFileInfo, error) {
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
	defer func() {
		_ = concatFile.Close()
	}()

	var segments []*SegmentInfo
	for _, info := range track.Infos {
		segments = append(segments, info.Segments...)
	}

	targetCodec := p.targetedCodec(track)

	// Add 1ms 1x1 black frame to avoid issue with switching resolution
	if track.TrackKind == trackKindVideo {
		blackFramePath := p.buildBlack1x1(config.OutputDir, targetCodec)
		if e := runFFmpegCommand(generateBlackVideoArguments(blackFramePath, targetCodec, 0.001, 1, 1, 1000), p.logger); e != nil {
			return nil, fmt.Errorf("failed to generate 1x1 black frame: %w", e)
		}
		cleanupFiles = append(cleanupFiles, blackFramePath)

		absPath, e := filepath.Abs(blackFramePath)
		if e != nil {
			return nil, e
		}

		if _, e := fmt.Fprintf(concatFile, "file '%s'\n", absPath); e != nil {
			return nil, e
		}
	}

	for i, segment := range segments {
		if track.TrackKind == trackKindVideo && !strings.EqualFold(segment.metadata.Codec, targetCodec) {
			if e := p.transcodeMinorityCodecFiles(targetCodec, segment); e != nil {
				return nil, e
			}
		}

		if _, e := fmt.Fprintf(concatFile, "file '%s'\n", segment.ContainerPath); e != nil {
			return nil, e
		}
		cleanupFiles = append(cleanupFiles, segment.ContainerPath)

		// Add gap filler if requested and there's a gap before the next segment
		if config.FillGap && i < len(segments)-1 {
			nextSegment := segments[i+1]
			offset := int64(0)
			if nextSegment.metadata.FirstKeyFrameOffsetMs != nil {
				offset = *nextSegment.metadata.FirstKeyFrameOffsetMs
			}
			gapDuration := offset + firstPacketNtpTimestamp(nextSegment.metadata) - segmentEndNtpTimestamp(segment.metadata)

			if gapDuration > 0 { // There's a gap
				gapSeconds := float64(gapDuration) / 1000.0
				p.logger.Infof("Detected %dms gap between segments, generating %s filler", gapDuration, track.TrackKind)

				// Create gap filler file
				gapFilePath := p.buildGapFilename(config.OutputDir, track, segment, i)

				var args []string
				if track.TrackKind == trackKindVideo {
					args = generateBlackVideoArguments(gapFilePath, segment.metadata.Codec, gapSeconds, 1280, 720, blackVideoFps)
				} else {
					args = generateSilenceArguments(gapFilePath, segment.ContainerExt, gapSeconds)
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

				if _, e := fmt.Fprintf(concatFile, "file '%s'\n", absPath); e != nil {
					return nil, e
				}
			}
		}
	}

	// Create final output file
	finalPath := p.buildFilename(config.OutputDir, track)

	// Concatenate all segments (with gap fillers if any)
	err = runFFmpegCommand(generateConcatFileArguments(finalPath, concatFile.Name()), p.logger)
	if err != nil {
		return nil, fmt.Errorf("failed to concatenate segments: %w", err)
	}

	p.logger.Debugf("Successfully concatenated %d segments into %s (gap filled %t)", len(segments), finalPath, config.FillGap)

	var ts, te int64
	if len(segments) > 0 {
		ts = segments[0].metadata.FirstRtpUnixTimestamp
		te = segments[len(segments)-1].metadata.LastRtpUnixTimestamp
	}

	var audioTrack, videoTrack []*TrackInfo
	switch track.TrackKind {
	case trackKindAudio:
		audioTrack = track.Infos
	case trackKindVideo:
		videoTrack = track.Infos
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

// transcodeMinorityCodecFiles re-encodes extracted MKV files whose codec differs
// from the target codec (the one with the longest total duration).
func (p *TrackExtractor) transcodeMinorityCodecFiles(targetCodec string, segment *SegmentInfo) error {
	p.logger.Infof("Transcoding %s from %s to %s", filepath.Base(segment.ContainerPath), segment.metadata.Codec, targetCodec)

	ext := "." + segment.ContainerExt
	transcodedPath := strings.TrimSuffix(segment.ContainerPath, ext) + "_transcoded" + ext
	if err := runFFmpegCommand(generateTranscodeVideoArguments(segment.ContainerPath, transcodedPath, targetCodec), p.logger); err != nil {
		return fmt.Errorf("failed to transcode %s: %w", filepath.Base(segment.ContainerPath), err)
	}

	segment.ContainerPath = transcodedPath
	segment.metadata.Codec = targetCodec
	return nil
}

// codecPreference defines transcoding cost order: lower value = cheaper to encode = preferred target.
var codecPreference = map[string]int{
	"video/h264": 0,
	"video/vp8":  1,
	"video/vp9":  2,
	"video/av1":  3,
}

// targetedCodec returns the preferred codec across all segments in the group.
// Preference order: H264 > VP8 > VP9 > AV1 (cheapest to encode first).
func (p *TrackExtractor) targetedCodec(group *aggregationGroup) string {
	bestCodec := ""
	bestPref := len(codecPreference) // worse than any known codec
	for _, info := range group.Infos {
		for _, s := range info.Segments {
			codec := strings.ToLower(s.metadata.Codec)
			pref, ok := codecPreference[codec]
			if !ok {
				continue
			} else if pref < bestPref {
				bestPref = pref
				bestCodec = codec
			}
		}
	}

	if bestCodec == "" {
		bestCodec = "video/h264"
	}
	return bestCodec
}

func (p *TrackExtractor) getMaxFrameDimension(track *aggregationGroup) SegmentFrameDimension {
	frameDimension := SegmentFrameDimension{}
	if track.TrackKind == trackKindVideo {
		for _, info := range track.Infos {
			for _, segment := range info.Segments {
				if segment.metadata.MaxFrameDimension != nil {
					frameDimension = getMaxFrameDimension(*segment.metadata.MaxFrameDimension, frameDimension)
				}
			}
		}
	}
	return frameDimension
}

// buildBlack1x1 creates an output filename for black 1x1 file
func (p *TrackExtractor) buildBlack1x1(outputDir string, codec string) string {
	return filepath.Join(outputDir, fmt.Sprintf("black_1x1_%s.%s", strings.ReplaceAll(codec, "/", "_"), mkvExtension))
}

// buildDefaultFilename creates output filename that indicates media type
func (p *TrackExtractor) buildFilename(outputDir string, group *aggregationGroup) string {
	track := group.Infos[0]
	return filepath.Join(outputDir, fmt.Sprintf("individual_%s_%s_%s_%d.%s", track.CallType, track.CallID, p.stringForGroup(group), track.CallStartTime.UnixMilli(), track.Segments[0].ContainerExt))
}

func (p *TrackExtractor) buildGapFilename(outputDir string, group *aggregationGroup, segment *SegmentInfo, i int) string {
	return filepath.Join(outputDir, fmt.Sprintf("gap_%s_%d_%d.%s", p.stringForGroup(group), group.Infos[0].CallStartTime.UnixMilli(), i, segment.ContainerExt))
}

func (p *TrackExtractor) buildConcatFilename(outputDir string, group *aggregationGroup) string {
	return filepath.Join(outputDir, fmt.Sprintf("concat_%s_%d.txt", p.stringForGroup(group), group.Infos[0].CallStartTime.UnixMilli()))
}

func (p *TrackExtractor) stringForGroup(group *aggregationGroup) string {
	userID := group.UserID
	sessionID := group.SessionID
	if sessionID != "" {
		sessionID = "_" + sessionID
	}
	trackID := group.TrackID
	if trackID != "" {
		trackID = "_" + trackID
	}
	media := group.TrackKind + "_only"
	if group.IsScreenshare {
		media = "shared_" + media
	}
	return fmt.Sprintf("%s%s%s_%s", userID, sessionID, trackID, media)
}
