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

	"github.com/GetStream/getstream-go/v3"
)

// TrackInfo represents a single track with its metadata (deduplicated across segments)
type TrackInfo struct {
	UserID        string         `json:"userId"`        // participant_id from timing metadata
	SessionID     string         `json:"sessionId"`     // user_session_id from timing metadata
	TrackID       string         `json:"trackId"`       // track_id from segment
	TrackType     string         `json:"trackType"`     // "audio" or "video" (cleaned from TRACK_TYPE_*)
	IsScreenshare bool           `json:"isScreenshare"` // true if this is a screenshare track
	Codec         string         `json:"codec"`         // codec info
	SegmentCount  int            `json:"segmentCount"`  // number of segments for this track
	Segments      []*SegmentInfo `json:"segments"`      // list of filenames (for JSON output only)

	ConcatenatedContainerPath string
}

type SegmentInfo struct {
	metadata *SegmentMetadata

	RtpDumpPath   string
	SdpPath       string
	ContainerPath string
	ContainerExt  string
	FFMpegOffset  int64
}

// RecordingMetadata contains all tracks and session information
type RecordingMetadata struct {
	Tracks   []*TrackInfo `json:"tracks"`
	UserIDs  []string     `json:"userIds"`
	Sessions []string     `json:"sessions"`
}

// MetadataParser handles parsing of raw recording files
type MetadataParser struct {
	logger *getstream.DefaultLogger
}

// NewMetadataParser creates a new metadata parser
func NewMetadataParser(logger *getstream.DefaultLogger) *MetadataParser {
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

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), "_timing_metadata.json") {
			p.logger.Debug("Processing metadata file: %s", path)

			data, err := os.ReadFile(path)
			if err != nil {
				p.logger.Warn("Failed to read metadata file %s: %v", path, err)
				return nil
			}

			tracks, err := p.parseTimingMetadataFile(data)
			if err != nil {
				p.logger.Warn("Failed to parse metadata file %s: %v", path, err)
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
	p.logger.Debug("Reading metadata directly from tar.gz (efficient mode): %s", tarGzPath)

	file, err := os.Open(tarGzPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tar.gz file: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

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

		// Only process timing metadata JSON files (skip all .rtpdump/.sdp files)
		if strings.HasSuffix(strings.ToLower(header.Name), "_timing_metadata.json") {
			p.logger.Debug("Processing metadata file: %s", header.Name)

			data, err := io.ReadAll(tarReader)
			if err != nil {
				p.logger.Warn("Failed to read metadata file %s: %v", header.Name, err)
				continue
			}

			tracks, err := p.parseTimingMetadataFile(data)
			if err != nil {
				p.logger.Warn("Failed to parse metadata file %s: %v", header.Name, err)
				continue
			}

			metadata.Tracks = append(metadata.Tracks, tracks...)
			filesRead++
		}
		// Skip all other files (.rtpdump, .sdp, etc.) - huge efficiency gain!
	}

	p.logger.Debug("Efficiently read %d metadata files from archive (skipped all media data files)", filesRead)

	// Extract unique user IDs and sessions
	metadata.UserIDs = p.extractUniqueUserIDs(metadata.Tracks)
	metadata.Sessions = p.extractUniqueSessions(metadata.Tracks)

	return metadata, nil
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

	processSegment := func(segment *SegmentMetadata, trackType string) {
		key := fmt.Sprintf("%s|%s|%s|%s",
			sessionMetadata.ParticipantID,
			sessionMetadata.UserSessionID,
			segment.TrackID,
			trackType)

		if existingTrack, exists := trackMap[key]; exists {
			existingTrack.Segments = append(existingTrack.Segments, &SegmentInfo{metadata: segment})
			existingTrack.SegmentCount++
		} else {
			// Create new track
			track := &TrackInfo{
				UserID:        sessionMetadata.ParticipantID,
				SessionID:     sessionMetadata.UserSessionID,
				TrackID:       segment.TrackID,
				TrackType:     p.cleanTrackType(segment.TrackType),
				IsScreenshare: p.isScreenshareTrack(segment.TrackType),
				Codec:         segment.Codec,
				SegmentCount:  1,
				Segments:      []*SegmentInfo{{metadata: segment}},
			}
			trackMap[key] = track
		}
	}

	// Process audio segments
	for _, segment := range sessionMetadata.Segments.Audio {
		processSegment(segment, p.cleanTrackType(segment.TrackType))
	}

	// Process video segments
	for _, segment := range sessionMetadata.Segments.Video {
		processSegment(segment, p.cleanTrackType(segment.TrackType))
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
		return "audio"
	case "TRACK_TYPE_VIDEO", "TRACK_TYPE_SCREEN_SHARE":
		return "video"
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
func FilterTracks(tracks []*TrackInfo, userID, sessionID, trackID, trackType, mediaFilter string) []*TrackInfo {
	filtered := make([]*TrackInfo, 0)

	for _, track := range tracks {
		if trackType != "" && track.TrackType != trackType {
			continue // Skip tracks with wrong TrackType
		}

		// Apply media type filtering if specified
		if mediaFilter != "" && mediaFilter != "both" {
			if mediaFilter == "user" && track.IsScreenshare {
				continue // Skip display tracks when only user requested
			}
			if mediaFilter == "display" && !track.IsScreenshare {
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

func firstPacketNtpTimestamp(segment *SegmentMetadata) int64 {
	if segment.FirstRtcpNtpTimestamp != 0 && segment.FirstRtcpRtpTimestamp != 0 {
		rtpNtpTs := (segment.FirstRtcpRtpTimestamp - segment.FirstRtpRtpTimestamp) / sampleRate(segment)
		return segment.FirstRtcpNtpTimestamp - int64(rtpNtpTs)
	} else {
		return segment.FirstRtpUnixTimestamp
	}
}

func lastPacketNtpTimestamp(segment *SegmentMetadata) int64 {
	if segment.LastRtcpNtpTimestamp != 0 && segment.LastRtcpRtpTimestamp != 0 {
		rtpNtpTs := (segment.LastRtpRtpTimestamp - segment.LastRtcpRtpTimestamp) / sampleRate(segment)
		return segment.LastRtcpNtpTimestamp + int64(rtpNtpTs)
	} else {
		return segment.LastRtpUnixTimestamp
	}
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
