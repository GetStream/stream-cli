package processing

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pion/rtp"
)

type GstreamerConverter struct {
	logger       *ProcessingLogger
	outputPath   string
	rtpConn      net.Conn
	gstreamerCmd *exec.Cmd
	mu           sync.Mutex
	ctx          context.Context
	cancel       context.CancelFunc
	port         int
	startAt      time.Time
}

func NewGstreamerConverter(outputPath, sdpContent string, logger *ProcessingLogger) (*GstreamerConverter, error) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &GstreamerConverter{
		logger:     logger,
		outputPath: outputPath,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Choose TCP listen port for GStreamer tcpserversrc
	r.port = rand.Intn(10000) + 10000

	// Start GStreamer with codec detection
	if err := r.startGStreamer(sdpContent, outputPath); err != nil {
		cancel()
		return nil, err
	}

	// Establish TCP client connection to the local tcpserversrc
	if err := r.setupConnections(r.port); err != nil {
		cancel()
		return nil, err
	}

	return r, nil
}

func (r *GstreamerConverter) setupConnections(port int) error {
	// Setup TCP connection with retry to match GStreamer tcpserversrc readiness
	address := "127.0.0.1:" + strconv.Itoa(port)
	deadline := time.Now().Add(10 * time.Second)
	var conn net.Conn
	var err error
	for {
		conn, err = net.DialTimeout("tcp", address, 500*time.Millisecond)
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("failed to connect to tcpserversrc at %s: %w", address, err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	r.rtpConn = conn
	return nil
}

func (r *GstreamerConverter) startGStreamer(sdpContent, outputFilePath string) error {
	r.startAt = time.Now()

	// Start with common GStreamer arguments optimized for RTP dump replay
	args, err := r.generateArgs(sdpContent, outputFilePath)
	if err != nil {
		return err
	}

	r.gstreamerCmd = exec.Command("gst-launch-1.0", args...)
	// Redirect output for debugging
	r.gstreamerCmd.Stdout = os.Stdout
	r.gstreamerCmd.Stderr = os.Stderr

	// Start GStreamer process
	if err := r.gstreamerCmd.Start(); err != nil {
		return err
	}

	r.logger.Infof("GStreamer process pid<%d> with pipeline: %s", r.gstreamerCmd.Process.Pid, strings.Join(args, " "))

	return nil
}

func (r *GstreamerConverter) generateArgs(sdpContent, outputFilePath string) ([]string, error) {
	// Parse SDP to determine RTP caps for rtpstreamdepay
	media, encodingName, payloadType, clockRate, err := parseRtpCapsFromSDP(sdpContent)
	if err != nil {
		return nil, err
	}

	// Start with common GStreamer arguments optimized for RTP dump replay
	args := []string{}
	args = append(args, "-e")
	// args = append(args, "--gst-debug-level=3")
	// args = append(args, "--gst-debug=tcpserversrc:5,rtp*:5,webm*:5,identity:5,jitterbuffer:5,av1*:5")
	// args = append(args, "--gst-debug-no-color")
	args = append(args, "tcpserversrc", "host=127.0.0.1", fmt.Sprintf("port=%d", r.port), "!")
	args = append(args, "application/x-rtp-stream", "!")
	args = append(args, "rtpstreamdepay", "!")
	args = append(args, fmt.Sprintf("application/x-rtp,media=%s,encoding-name=%s,clock-rate=%s,payload=%s", media, encodingName, clockRate, payloadType), "!")

	// Simplified approach for RTP dump replay:
	// - rtpjitterbuffer: Basic packet reordering with minimal interference
	//   - mode=none: Don't override timing, let depayloaders handle it
	//   - latency=0: No artificial latency, process packets as they come
	//   - do-retransmission=false: No retransmission for dump replay
	args = append(args, "rtpjitterbuffer", "mode=none", "latency=0", "do-lost=false", "do-retransmission=false", "drop-on-latency=false", "!")

	switch encodingName {
	case "VP9", "AV1", "H264":
		args = append(args, fmt.Sprintf("rtp%sdepay", strings.ToLower(encodingName)), "!")
		args = append(args, fmt.Sprintf("%sparse", strings.ToLower(encodingName)), "!")
	case "OPUS", "VP8":
		args = append(args, fmt.Sprintf("rtp%sdepay", strings.ToLower(encodingName)), "!")
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encodingName)
	}

	args = append(args, "matroskamux", "streamable=false", "!")
	args = append(args, "filesink", fmt.Sprintf("location=%s", outputFilePath))

	return args, nil
}

func parseRtpCapsFromSDP(sdp string) (media string, encodingName string, payload string, clockRate string, err error) {
	// Expect one m= line and one a=rtpmap line; return error if missing or malformed
	mLineFound := false
	rtpmapLineFound := false
	for _, raw := range strings.Split(sdp, "\n") {
		//line := strings.TrimSpace(raw)
		lower := strings.ToLower(raw)
		if strings.HasPrefix(lower, "m=") {
			mLineFound = true
			// Format: m=<media> <port> <proto> <pt> <pt> ...
			fields := strings.Fields(lower)
			if len(fields) >= 1 {
				media = strings.TrimPrefix(fields[0], "m=")
			} else {
				err = fmt.Errorf("invalid m= line: %s", lower)
				return
			}
		} else if strings.HasPrefix(lower, "a=rtpmap:") {
			rtpmapLineFound = true

			// Format: a=rtpmap:<pt> <encoding>/<clock>[/channels]
			after := strings.TrimSpace(lower[len("a=rtpmap:"):])
			fields := strings.Fields(after)
			if len(fields) >= 2 {
				payload = fields[0]
				codec := strings.ToUpper(fields[1])
				parts := strings.Split(codec, "/")
				if len(parts) >= 2 {
					encodingName = parts[0]
					clockRate = parts[1]
				} else {
					err = fmt.Errorf("invalid a=rtpmap: %s", lower)
					return
				}
			} else {
				err = fmt.Errorf("invalid a=rtpmap: %s", lower)
				return
			}
		}
	}

	if !mLineFound || !rtpmapLineFound {
		err = fmt.Errorf("Invalid SDP m= or a=rtpmap lines not found: \n%s", sdp)
	}
	return
}

func (r *GstreamerConverter) OnRTP(packet *rtp.Packet) error {
	// Marshal RTP packet
	buf, err := packet.Marshal()
	if err != nil {
		return err
	}

	return r.PushRtpBuf(buf)
}

func (r *GstreamerConverter) PushRtpBuf(buf []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Send RTP packet over TCP using RFC4571 2-byte length prefix
	if r.rtpConn != nil {
		if len(buf) > 0xFFFF {
			return fmt.Errorf("rtp packet too large for TCP framing: %d bytes", len(buf))
		}
		header := make([]byte, 2)
		binary.BigEndian.PutUint16(header, uint16(len(buf)))
		if _, err := r.rtpConn.Write(header); err != nil {
			r.logger.Warnf("Failed to write RTP length header: %v", err)
			return err
		}
		if _, err := r.rtpConn.Write(buf); err != nil {
			r.logger.Warnf("Failed to write RTP packet: %v", err)
			return err
		}
	}
	return nil
}

func (r *GstreamerConverter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Infof("GStreamer process pid<%d> Closing TCP connection and wait for termination...", r.gstreamerCmd.Process.Pid)

	// Cancel context to stop background goroutines
	if r.cancel != nil {
		r.cancel()
	}

	// Close TCP connection
	if r.rtpConn != nil {
		_ = r.rtpConn.Close()
		r.rtpConn = nil
	}

	// Gracefully wait for FFmpeg termination
	if r.gstreamerCmd != nil && r.gstreamerCmd.Process != nil {
		// Wait for graceful exit with timeout
		done := make(chan error, 1)
		go func() {
			done <- r.gstreamerCmd.Wait()
		}()

		select {
		case <-time.After(5 * time.Second):
			r.logger.Warnf("GStreamer process pid<%d> termination timeout in %s...", r.gstreamerCmd.Process.Pid, time.Since(r.startAt).Round(time.Millisecond))

			// Timeout, force kill
			if e := r.gstreamerCmd.Process.Kill(); e != nil {
				r.logger.Errorf("GStreamer process pid<%d> errored while killing: %v", r.gstreamerCmd.Process.Pid, e)
			}
		case <-done:
			r.logger.Infof("GStreamer process pid<%d> exited succesfully in %s...", r.gstreamerCmd.Process.Pid, time.Since(r.startAt).Round(time.Millisecond))
		}
	}

	return nil
}
