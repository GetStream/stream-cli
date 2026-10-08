package rawrecording

// Flag names for global/persistent flags
const (
	FlagInputFile = "input-file"
	FlagInputDir  = "input-dir"
	FlagInputS3   = "input-s3"
	FlagOutput    = "output"
	FlagVerbose   = "verbose"
	FlagCacheDir  = "cache-dir"
)

// Flag names for filter flags (used across multiple commands)
const (
	FlagUserID    = "user-id"
	FlagSessionID = "session-id"
	FlagTrackID   = "track-id"
)

// Flag names for processing options
const (
	FlagFillGaps = "fill-gaps"
	FlagFixDtx   = "fix-dtx"
	FlagMedia    = "media"
)

// Flag names for list-tracks command
const (
	FlagFormat         = "format"
	FlagTrackType      = "track-type"
	FlagCompletionType = "completion-type"
)

// Flag descriptions for global/persistent flags
const (
	DescInputFile = "Raw recording tar.gz archive path"
	DescInputDir  = "Raw recording extracted directory path"
	DescInputS3   = "Raw recording S3 URL (s3://bucket/path or presigned HTTPS URL)"
	DescOutput    = "Output directory"
	DescVerbose   = "Enable verbose logging"
	DescCacheDir  = "Cache directory for S3 downloads"
)

// Flag descriptions for filter flags
const (
	DescUserID    = "Filter by user ID"
	DescSessionID = "Filter by session ID"
	DescTrackID   = "Filter by track ID"
)

// Flag descriptions for processing options
const (
	DescFillGapsAudio = "Fill with silence when track was muted"
	DescFillGapsVideo = "Fill with black frame when track was muted"
	DescFixDtx        = "Restore original audio duration by filling DTX silence gaps (required for A/V sync)"
	DescMedia         = "Filter by media type: 'user', 'display', or 'both'"
)

// Flag descriptions for list-tracks command
const (
	DescFormat         = "Output format: table, json, users, sessions, tracks, completion"
	DescTrackType      = "Filter by track type: audio, video"
	DescCompletionType = "For completion format: users, sessions, tracks"
)

// Default values
const (
	DefaultFormat         = "table"
	DefaultCompletionType = "tracks"
	DefaultMedia          = "both"
	DefaultCacheSubdir    = "stream-cli/raw-recordings"
)

// Media type values
const (
	MediaUser    = "user"
	MediaDisplay = "display"
	MediaBoth    = "both"
)

// Track type values
const (
	TrackTypeAudio = "audio"
	TrackTypeVideo = "video"
)
