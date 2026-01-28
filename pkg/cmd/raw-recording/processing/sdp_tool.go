package processing

import (
	"fmt"
	"os"
	"strings"

	"github.com/pion/webrtc/v4"
)

func readSDP(sdpFilePath string) (string, error) {
	content, err := os.ReadFile(sdpFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read SDP file %s: %w", sdpFilePath, err)
	}
	return string(content), nil
}

func replaceSDP(sdpContent string, port int) string {
	lines := strings.Split(sdpContent, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "m=") {
			// Parse the m= line: m=<media_type> <port> RTP/AVP <payload_type>
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				// Replace the port (second field)
				parts[1] = fmt.Sprintf("%d", port)
				lines[i] = strings.Join(parts, " ")
				break
			}
		}
	}
	return strings.Join(lines, "\n")
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
