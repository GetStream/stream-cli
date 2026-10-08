package processing

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// TrackInfo represents a single track with its metadata (deduplicated across segments)
type TrackInfo struct {
	CallType       string         `json:"callType"`       // call_type from timing metadata
	CallID         string         `json:"callId"`         // call_id from timing metadata
	CallSessionID  string         `json:"callSessionId"`  // call_session_id from timing metadata
	CallStartTime  time.Time      `json:"callStartTime"`  // call_start_time from timing metadata
	UserID         string         `json:"userId"`         // participant_id from timing metadata
	SessionID      string         `json:"sessionId"`      // user_session_id from timing metadata
	TrackID        string         `json:"trackId"`        // track_id from segment
	TrackType      string         `json:"trackType"`      // track_type from segment
	TrackKind      string         `json:"trackKind"`      // "audio" or "video" (cleaned from TRACK_TYPE_*)
	IsScreenshare  bool           `json:"isScreenshare"`  // true if this is a screenshare track
	Codec          string         `json:"codec"`          // codec info
	SegmentCount   int            `json:"segmentCount"`   // number of segments for this track
	TrackStartTime time.Time      `json:"trackStartTime"` // first_rtp_unix_timestamp from segment
	TrackEndTime   time.Time      `json:"trackEndTime"`   // last_rtp_unix_timestamp from segment
	Segments       []*SegmentInfo `json:"segments"`       // list of filenames (for JSON output only)
}

type SegmentInfo struct {
	metadata *SegmentMetadata

	RtpDumpPath   string
	SdpPath       string
	ContainerPath string
	ContainerExt  string
}

type TrackFileInfo struct {
	Name              string
	StartAt           time.Time
	EndAt             time.Time
	MaxFrameDimension SegmentFrameDimension
	AudioTrack        []*TrackInfo
	VideoTrack        []*TrackInfo
}

// CallMetadataFilename is the canonical filename for call-level metadata,
// written once per recording at the root of the archive.
const CallMetadataFilename = "call_metadata.json"

// CallMetadata holds call-level information about a recording: identity and
// the recording window (start/end). It is written once per call at the root
// of the recording directory and is independent of per-participant timing
// files. Absent in archives produced before this file was introduced.
type CallMetadata struct {
	CallType           string    `json:"call_type"`
	CallID             string    `json:"call_id"`
	CallSessionID      string    `json:"call_session_id"`
	RecordingStartTime time.Time `json:"recording_start_time"`
	RecordingEndTime   time.Time `json:"recording_end_time,omitempty"`
}

// RecordingMetadata contains all tracks and session information
type RecordingMetadata struct {
	CallMetadata *CallMetadata `json:"callMetadata,omitempty"`
	Tracks       []*TrackInfo  `json:"tracks"`
	UserIDs      []string      `json:"userIds"`
	Sessions     []string      `json:"sessions"`
}

// MetadataParser handles parsing of raw recording files
type MetadataParser struct {
	logger *ProcessingLogger
}

// NewMetadataParser creates a new metadata parser
func NewMetadataParser(logger *ProcessingLogger) *MetadataParser {
	return &MetadataParser{
		logger: logger,
	}
}

// ParseMetadataOnly efficiently extracts only metadata from archives (optimized for list-tracks)
// This is much faster than full extraction when you only need timing metadata
func (p *MetadataParser) ParseMetadataOnly(inputPath string) (*RecordingMetadata, error) {
	// If it's already a directory, use the normal path
	if stat, err := os.Stat(inputPath); err == nil && stat.IsDir() {
		return p.parseDirectory(inputPath)
	}

	// If it's a tar.gz file, use selective extraction (much faster)
	if strings.HasSuffix(strings.ToLower(inputPath), ".tar.gz") {
		return p.parseMetadataOnlyFromTarGz(inputPath)
	}

	return nil, fmt.Errorf("unsupported input format: %s (only tar.gz files and directories supported)", inputPath)
}

// parseDirectory processes a directory containing recording files
func (p *MetadataParser) parseDirectory(dirPath string) (*RecordingMetadata, error) {
	metadata := &RecordingMetadata{
		Tracks:   make([]*TrackInfo, 0),
		UserIDs:  make([]string, 0),
		Sessions: make([]string, 0),
	}

	// Find and process timing metadata files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		name := strings.ToLower(info.Name())
		switch {
		case name == CallMetadataFilename:
			p.logger.Debugf("Processing call metadata file: %s", path)
			data, err := os.ReadFile(path)
			if err != nil {
				p.logger.Warnf("Failed to read call metadata file %s: %v", path, err)
				return nil
			}
			cm, err := p.parseCallMetadataFile(data)
			if err != nil {
				p.logger.Warnf("Failed to parse call metadata file %s: %v", path, err)
				return nil
			}
			metadata.CallMetadata = cm
		case strings.HasSuffix(name, "_timing_metadata.json"):
			p.logger.Debugf("Processing metadata file: %s", path)

			data, err := os.ReadFile(path)
			if err != nil {
				p.logger.Warnf("Failed to read metadata file %s: %v", path, err)
				return nil
			}

			tracks, err := p.parseTimingMetadataFile(data)
			if err != nil {
				p.logger.Warnf("Failed to parse metadata file %s: %v", path, err)
				return nil
			}

			metadata.Tracks = append(metadata.Tracks, tracks...)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to process directory: %w", err)
	}

	// Build unique lists
	metadata.UserIDs = p.extractUniqueUserIDs(metadata.Tracks)
	metadata.Sessions = p.extractUniqueSessions(metadata.Tracks)

	return metadata, nil
}

// parseMetadataOnlyFromTarGz efficiently extracts only timing metadata from tar.gz files
// This is optimized for list-tracks - only reads JSON files, skips all .rtpdump/.sdp files
func (p *MetadataParser) parseMetadataOnlyFromTarGz(tarGzPath string) (*RecordingMetadata, error) {
	p.logger.Debugf("Reading metadata directly from tar.gz (efficient mode): %s", tarGzPath)

	file, err := os.Open(tarGzPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tar.gz file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer func() {
		_ = gzReader.Close()
	}()

	tarReader := tar.NewReader(gzReader)

	metadata := &RecordingMetadata{
		Tracks:   make([]*TrackInfo, 0),
		UserIDs:  make([]string, 0),
		Sessions: make([]string, 0),
	}

	filesRead := 0
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("failed to read tar entry: %w", err)
		} else if header.FileInfo().IsDir() {
			continue
		}

		// Only process metadata JSON files (skip all .rtpdump/.sdp files)
		name := strings.ToLower(filepath.Base(header.Name))
		switch {
		case name == CallMetadataFilename:
			p.logger.Debugf("Processing call metadata file: %s", header.Name)
			data, err := io.ReadAll(tarReader)
			if err != nil {
				p.logger.Warnf("Failed to read call metadata file %s: %v", header.Name, err)
				continue
			}
			cm, err := p.parseCallMetadataFile(data)
			if err != nil {
				p.logger.Warnf("Failed to parse call metadata file %s: %v", header.Name, err)
				continue
			}
			metadata.CallMetadata = cm
			filesRead++
		case strings.HasSuffix(name, "_timing_metadata.json"):
			p.logger.Debugf("Processing metadata file: %s", header.Name)

			data, err := io.ReadAll(tarReader)
			if err != nil {
				p.logger.Warnf("Failed to read metadata file %s: %v", header.Name, err)
				continue
			}

			tracks, err := p.parseTimingMetadataFile(data)
			if err != nil {
				p.logger.Warnf("Failed to parse metadata file %s: %v", header.Name, err)
				continue
			}

			metadata.Tracks = append(metadata.Tracks, tracks...)
			filesRead++
		}
		// Skip all other files (.rtpdump, .sdp, etc.) - huge efficiency gain!
	}

	p.logger.Debugf("Efficiently read %d metadata files from archive (skipped all media data files)", filesRead)

	// Extract unique user IDs and sessions
	metadata.UserIDs = p.extractUniqueUserIDs(metadata.Tracks)
	metadata.Sessions = p.extractUniqueSessions(metadata.Tracks)

	return metadata, nil
}

// parseCallMetadataFile parses the call-level metadata JSON file
func (p *MetadataParser) parseCallMetadataFile(data []byte) (*CallMetadata, error) {
	var cm CallMetadata
	if err := json.Unmarshal(data, &cm); err != nil {
		return nil, fmt.Errorf("failed to unmarshal call metadata: %w", err)
	}
	return &cm, nil
}

// parseTimingMetadataFile parses a timing metadata JSON file and extracts tracks
func (p *MetadataParser) parseTimingMetadataFile(data []byte) ([]*TrackInfo, error) {
	var sessionMetadata SessionTimingMetadata
	err := json.Unmarshal(data, &sessionMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session metadata: %w", err)
	}

	// Use a map to deduplicate tracks by unique key
	trackMap := make(map[string]*TrackInfo)

	processSegment := func(segment *SegmentMetadata) {
		key := fmt.Sprintf("%s|%s|%s|%s",
			sessionMetadata.ParticipantID,
			sessionMetadata.UserSessionID,
			segment.TrackID,
			strings.ReplaceAll(strings.ToLower(segment.Codec), "/", "|"))

		if existingTrack, exists := trackMap[key]; exists {
			existingTrack.Segments = append(existingTrack.Segments, &SegmentInfo{metadata: segment})
			existingTrack.SegmentCount++

			ts, te := time.UnixMilli(segment.FirstRtpUnixTimestamp), time.UnixMilli(segment.LastRtpUnixTimestamp)
			if ts.Before(existingTrack.TrackStartTime) {
				existingTrack.TrackStartTime = ts
			}
			if te.After(existingTrack.TrackEndTime) {
				existingTrack.TrackEndTime = te
			}
		} else {
			// Create new track
			track := &TrackInfo{
				CallType:       sessionMetadata.CallType,
				CallID:         sessionMetadata.CallID,
				CallSessionID:  sessionMetadata.CallSessionID,
				CallStartTime:  sessionMetadata.CallStartTime,
				UserID:         sessionMetadata.ParticipantID,
				SessionID:      sessionMetadata.UserSessionID,
				TrackID:        segment.TrackID,
				TrackType:      segment.TrackType,
				TrackKind:      p.cleanTrackType(segment.TrackType),
				IsScreenshare:  p.isScreenshareTrack(segment.TrackType),
				Codec:          segment.Codec,
				SegmentCount:   1,
				TrackStartTime: time.UnixMilli(segment.FirstRtpUnixTimestamp),
				TrackEndTime:   time.UnixMilli(segment.LastRtpUnixTimestamp),
				Segments:       []*SegmentInfo{{metadata: segment}},
			}
			trackMap[key] = track
		}
	}

	// Process audio segments
	for _, segment := range sessionMetadata.Segments.Audio {
		processSegment(segment)
	}

	// Process video segments
	for _, segment := range sessionMetadata.Segments.Video {
		processSegment(segment)
	}

	// Convert map to slice
	tracks := make([]*TrackInfo, 0, len(trackMap))
	for _, track := range trackMap {
		sort.Slice(track.Segments, func(i, j int) bool {
			return track.Segments[i].metadata.FirstRtpUnixTimestamp < track.Segments[j].metadata.FirstRtpUnixTimestamp
		})
		tracks = append(tracks, track)
	}

	return tracks, nil
}

// isScreenshareTrack detects if a track is screenshare-related
func (p *MetadataParser) isScreenshareTrack(trackType string) bool {
	return trackType == "TRACK_TYPE_SCREEN_SHARE_AUDIO" || trackType == "TRACK_TYPE_SCREEN_SHARE"
}

// cleanTrackType converts TRACK_TYPE_* to simple "audio" or "video"
func (p *MetadataParser) cleanTrackType(trackType string) string {
	switch trackType {
	case "TRACK_TYPE_AUDIO", "TRACK_TYPE_SCREEN_SHARE_AUDIO":
		return trackKindAudio
	case "TRACK_TYPE_VIDEO", "TRACK_TYPE_SCREEN_SHARE":
		return trackKindVideo
	default:
		return strings.ToLower(trackType)
	}
}

// extractUniqueUserIDs returns a sorted list of unique user IDs
func (p *MetadataParser) extractUniqueUserIDs(tracks []*TrackInfo) []string {
	userIDMap := make(map[string]bool)
	for _, track := range tracks {
		userIDMap[track.UserID] = true
	}

	userIDs := make([]string, 0, len(userIDMap))
	for userID := range userIDMap {
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

// NOTE: ExtractTrackFiles and extractTrackFromTarGz removed - no longer needed since we always work with directories

// extractUniqueSessions returns a sorted list of unique session IDs
func (p *MetadataParser) extractUniqueSessions(tracks []*TrackInfo) []string {
	sessionMap := make(map[string]bool)
	for _, track := range tracks {
		sessionMap[track.SessionID] = true
	}

	sessions := make([]string, 0, len(sessionMap))
	for session := range sessionMap {
		sessions = append(sessions, session)
	}

	return sessions
}

// FilterTracks filters tracks based on mutually exclusive criteria
// Only one filter (userID, sessionID, or trackID) can be specified at a time
// Empty values are ignored, specific values must match
// If all are empty, all tracks are returned
func FilterTracks(tracks []*TrackInfo, userID, sessionID, trackID, trackKind, mediaType string) []*TrackInfo {
	filtered := make([]*TrackInfo, 0)

	for _, track := range tracks {
		if trackKind != "" && track.TrackKind != trackKind {
			continue // Skip tracks with wrong trackKind
		}

		// Apply media type filtering if specified
		if mediaType != "" && mediaType != MediaTypeBoth {
			if mediaType == MediaTypeUser && track.IsScreenshare {
				continue // Skip display tracks when only user requested
			}
			if mediaType == MediaTypeDisplay && !track.IsScreenshare {
				continue // Skip user tracks when only display requested
			}
		}

		// Apply the single specified filter (mutually exclusive)
		if trackID != "" {
			// Filter by trackID - return only that specific track
			if track.TrackID == trackID {
				filtered = append(filtered, track)
			}
		} else if sessionID != "" {
			// Filter by sessionID - return all tracks for that session
			if track.SessionID == sessionID {
				filtered = append(filtered, track)
			}
		} else if userID != "" {
			// Filter by userID - return all tracks for that user
			if track.UserID == userID {
				filtered = append(filtered, track)
			}
		} else {
			// No filters specified - return all tracks
			filtered = append(filtered, track)
		}
	}

	return filtered
}

func firstPacketTimestamp(segmentMetadata *SegmentMetadata) int64 {
	firstOffset := int64(0)
	if segmentMetadata.FirstKeyFrameOffsetMs != nil {
		firstOffset = *segmentMetadata.FirstKeyFrameOffsetMs
	}

	firstTs := firstOffset + firstPacketNtpTimestamp(segmentMetadata)
	return firstTs
}

func firstPacketNtpTimestamp(segment *SegmentMetadata) int64 {
	if segment.FirstRtcpNtpTimestamp != 0 && segment.FirstRtcpRtpTimestamp != 0 {
		return segment.FirstRtcpNtpTimestamp + rtpTimestampDeltaMs(segment.FirstRtpRtpTimestamp, segment.FirstRtcpRtpTimestamp, segment)
	}
	return segment.FirstRtpUnixTimestamp
}

func lastPacketNtpTimestamp(segment *SegmentMetadata) int64 {
	if segment.LastRtcpNtpTimestamp != 0 && segment.LastRtcpRtpTimestamp != 0 {
		return segment.LastRtcpNtpTimestamp + rtpTimestampDeltaMs(segment.LastRtpRtpTimestamp, segment.LastRtcpRtpTimestamp, segment)
	}
	return segment.LastRtpUnixTimestamp
}

// segmentEndNtpTimestamp returns the wall time of a segment's last packet, anchored on
// the first RTCP SR and spanned by the segment's own RTP clock.
func segmentEndNtpTimestamp(segment *SegmentMetadata) int64 {
	return firstPacketNtpTimestamp(segment) + segmentRtpDurationMs(segment)
}

// segmentRtpDurationMs returns the segment's span on the sender's RTP clock, preferring
// the recorder's own RtpDurationMs so the two sides cannot disagree. Archives written
// before that field existed omit it, so fall back to deriving the span from the segment's
// RTP timestamps.
func segmentRtpDurationMs(segment *SegmentMetadata) int64 {
	if segment.RtpDurationMs > 0 {
		return segment.RtpDurationMs
	}
	if segment.LastRtpRtpTimestamp == 0 {
		return 0
	}
	// Mirrors the recorder: a negative delta means no forward progress, not a wrap.
	if d := rtpTimestampDeltaMs(segment.LastRtpRtpTimestamp, segment.FirstRtpRtpTimestamp, segment); d > 0 {
		return d
	}
	return 0
}

// rtpTimestampDeltaMs returns the signed millisecond delta (rtp - ref) between two
// 32-bit RTP timestamps, added to an RTCP NTP reference to recover a packet's wall time.
//
// RTP timestamps are uint32 and wrap at 2^32. A naive numeric comparison mis-signs the
// delta whenever the two values straddle the wrap boundary — e.g. an RTCP SR reporting
// rtp_timestamp=123582 that is actually *later* than an RTP packet at 4294614425. Casting
// the modular difference to int32 yields the signed shortest distance, which is correct as
// long as the true gap stays under 2^31 samples (~12.4h at 48kHz, ~6.6h at 90kHz).
func rtpTimestampDeltaMs(rtp, ref uint32, segment *SegmentMetadata) int64 {
	return int64(int32(rtp-ref)) / int64(sampleRate(segment))
}

func sampleRate(segment *SegmentMetadata) uint32 {
	switch segment.TrackType {
	case "TRACK_TYPE_AUDIO",
		"TRACK_TYPE_SCREEN_SHARE_AUDIO":
		return 48
	default:
		return 90
	}
}
