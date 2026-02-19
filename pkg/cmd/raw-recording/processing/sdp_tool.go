package processing

import (
	"fmt"
	"os"
	"strings"

	webrtc "github.com/pion/webrtc/v4"
)

func readSDP(sdpFilePath string) (string, error) {
	content, err := os.ReadFile(sdpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read SDP file %s: %w", sdpFilePath, err)
	}
	return string(content), nil
}

func mimeType(sdp string) (string, error) {
	upper := strings.ToUpper(sdp)
	if strings.Contains(upper, "VP9") {
		return webrtc.MimeTypeVP9, nil
	}
	if strings.Contains(upper, "VP8") {
		return webrtc.MimeTypeVP8, nil
	}
	if strings.Contains(upper, "AV1") {
		return webrtc.MimeTypeAV1, nil
	}
	if strings.Contains(upper, "OPUS") {
		return webrtc.MimeTypeOpus, nil
	}
	if strings.Contains(upper, "H264") {
		return webrtc.MimeTypeH264, nil
	}

	return "", fmt.Errorf("mimeType should be OPUS, VP8, VP9, AV1, H264")
}
