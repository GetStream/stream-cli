package processing

type SessionTimingMetadata struct {
	ParticipantID string `json:"participant_id"`
	UserSessionID string `json:"user_session_id"`
	Segments      struct {
		Audio []*SegmentMetadata `json:"audio"`
		Video []*SegmentMetadata `json:"video"`
	} `json:"segments"`
}

type SegmentMetadata struct {
	// Global information
	BaseFilename string `json:"base_filename"`

	// Track information
	Codec     string `json:"codec"`
	TrackID   string `json:"track_id"`
	TrackType string `json:"track_type"`

	// Packet timing information
	FirstRtpRtpTimestamp  uint32 `json:"first_rtp_rtp_timestamp"`
	FirstRtpUnixTimestamp int64  `json:"first_rtp_unix_timestamp"`
	LastRtpRtpTimestamp   uint32 `json:"last_rtp_rtp_timestamp,omitempty"`
	LastRtpUnixTimestamp  int64  `json:"last_rtp_unix_timestamp,omitempty"`
	FirstRtcpRtpTimestamp uint32 `json:"first_rtcp_rtp_timestamp,omitempty"`
	FirstRtcpNtpTimestamp int64  `json:"first_rtcp_ntp_timestamp,omitempty"`
	LastRtcpRtpTimestamp  uint32 `json:"last_rtcp_rtp_timestamp,omitempty"`
	LastRtcpNtpTimestamp  int64  `json:"last_rtcp_ntp_timestamp,omitempty"`
}
