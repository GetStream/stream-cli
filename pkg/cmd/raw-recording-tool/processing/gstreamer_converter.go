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

	"github.com/GetStream/getstream-go/v3"
	"github.com/pion/rtp"
)

type CursorGstreamerWebmRecorder struct {
	logger          *getstream.DefaultLogger
	outputPath      string
	rtpConn         net.Conn
	gstreamerCmd    *exec.Cmd
	mu              sync.Mutex
	ctx             context.Context
	cancel          context.CancelFunc
	port            int
	sdpFile         *os.File
	finalOutputPath string // Path for post-processed file with duration
	tempOutputPath  string // Path for temporary file before post-processing
}

func NewCursorGstreamerWebmRecorder(outputPath, sdpContent string, logger *getstream.DefaultLogger) (*CursorGstreamerWebmRecorder, error) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &CursorGstreamerWebmRecorder{
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

func (r *CursorGstreamerWebmRecorder) setupConnections(port int) error {
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
		time.Sleep(100 * time.Millisecond)
	}
	r.rtpConn = conn
	return nil
}

func (r *CursorGstreamerWebmRecorder) startGStreamer(sdpContent, outputFilePath string) error {
	// Parse SDP to determine RTP caps for rtpstreamdepay
	media, encodingName, payloadType, clockRate := parseRtpCapsFromSDP(sdpContent)
	r.logger.Info("Starting TCP-based GStreamer pipeline (media=%s, encoding=%s, payload=%d, clock-rate=%d)", media, encodingName, payloadType, clockRate)

	// Determine codec from SDP content and build GStreamer arguments
	isVP9 := strings.Contains(strings.ToUpper(sdpContent), "VP9")
	isVP8 := strings.Contains(strings.ToUpper(sdpContent), "VP8")
	isAV1 := strings.Contains(strings.ToUpper(sdpContent), "AV1")
	isH264 := strings.Contains(strings.ToUpper(sdpContent), "H264") || strings.Contains(strings.ToUpper(sdpContent), "H.264")
	isOpus := strings.Contains(strings.ToUpper(sdpContent), "OPUS")

	// Start with common GStreamer arguments optimized for RTP dump replay
	args := []string{
		//"--gst-debug-level=3",
		//"--gst-debug=tcpserversrc:5,rtp*:5,webm*:5,identity:5,jitterbuffer:5,vp9*:5",
		//"--gst-debug-no-color",
		"-e", // Send EOS on interrupt for clean shutdown
	}
	// Source from TCP (RFC4571 framed) and depayload back to application/x-rtp
	args = append(args,
		"tcpserversrc",
		"host=127.0.0.1",
		fmt.Sprintf("port=%d", r.port),
		"name=tcp_in",
		"!",
		"queue",
		"max-size-buffers=0",
		"max-size-bytes=268435456",
		"max-size-time=0",
		"leaky=0",
		"!",
		// Ensure rtpstreamdepay sink has caps
		"application/x-rtp-stream",
		"!",
		"rtpstreamdepay",
		"!",
		fmt.Sprintf("application/x-rtp,media=%s,encoding-name=%s,clock-rate=%d,payload=%d", media, encodingName, clockRate, payloadType),
		"!",
	)

	// Build pipeline based on codec with simplified RTP timestamp handling for dump replay
	//
	// Simplified approach for RTP dump replay:
	// - rtpjitterbuffer: Basic packet reordering with minimal interference
	//   - latency=0: No artificial latency, process packets as they come
	//   - mode=none: Don't override timing, let depayloaders handle it
	//   - do-retransmission=false: No retransmission for dump replay
	// - Remove identity sync to avoid timing conflicts
	//
	// This approach focuses on preserving original RTP timestamps without
	// artificial buffering that can interfere with dump replay timing.
	if false && isH264 {
		r.logger.Info("Detected H.264 codec, building H.264 pipeline with timestamp handling...")
		args = append(args,
			"application/x-rtp,media=video,encoding-name=H264,clock-rate=90000", "!",
			"rtpjitterbuffer",
			"latency=0",
			"mode=none",
			"do-retransmission=false", "!",
			"rtph264depay", "!",
			"h264parse", "!",
			"mp4mux", "!",
			"filesink", fmt.Sprintf("location=%s", outputFilePath),
		)
	} else if false && isVP9 {
		r.logger.Info("Detected VP9 codec, building VP9 pipeline with timestamp handling...")
		args = append(args,
			"rtpjitterbuffer",
			"latency=0",
			"mode=none",
			"do-retransmission=false",
			"drop-on-latency=false",
			"buffer-mode=slave",
			"max-dropout-time=5000000000",
			"max-reorder-delay=1000000000",
			"!",
			"rtpvp9depay", "!",
			"vp9parse", "!",
			"webmmux",
			"writing-app=GStreamer-VP9",
			"streamable=false",
			"min-index-interval=2000000000", "!",
			"filesink", fmt.Sprintf("location=%s", outputFilePath),
		)
	} else if isVP9 {
		r.logger.Info("Detected VP9 codec, building VP9 pipeline with RTP timestamp handling...")
		args = append(args,

			//// jitterbuffer for packet reordering and timestamp handling
			"rtpjitterbuffer",
			"name=jitterbuffer",
			"mode=none",
			"latency=0",               // No artificial latency - process immediately
			"do-lost=false",           // Don't generate lost events for missing packets
			"do-retransmission=false", // No retransmission for offline replay
			"drop-on-latency=false",   // Keep all packets even if late
			"!",
			//
			// Depayload RTP to get VP9 frames
			"rtpvp9depay",
			"!",

			// Parse VP9 stream to ensure valid frame structure
			"vp9parse",
			"!",

			// Queue for buffering
			"queue",
			"!",

			// Mux into Matroska/WebM container
			"webmmux",
			"writing-app=GStreamer-VP9",
			"streamable=false",
			"min-index-interval=2000000000",
			"!",

			// Write to file
			"filesink",
			fmt.Sprintf("location=%s", outputFilePath),
		)

	} else if false && isVP8 {
		r.logger.Info("Detected VP8 codec, building VP8 pipeline with timestamp handling...")
		args = append(args,
			"application/x-rtp,media=video,encoding-name=VP8,clock-rate=90000", "!",
			"rtpjitterbuffer",
			"latency=0",
			"mode=none",
			"do-retransmission=false", "!",
			"rtpvp8depay", "!",
			"vp8parse", "!",
			"webmmux", "writing-app=GStreamer", "streamable=false", "min-index-interval=2000000000", "!",
			"filesink", fmt.Sprintf("location=%s", outputFilePath),
		)
	} else if false && isAV1 {
		r.logger.Info("Detected AV1 codec, building AV1 pipeline with timestamp handling...")
		args = append(args,
			"application/x-rtp,media=video,encoding-name=AV1,clock-rate=90000", "!",
			"rtpjitterbuffer",
			"latency=0",
			"mode=none",
			"do-retransmission=false", "!",
			"rtpav1depay", "!",
			"av1parse", "!",
			"webmmux", "!",
			"filesink", fmt.Sprintf("location=%s", outputFilePath),
		)
	} else if false && isOpus {
		r.logger.Info("Detected Opus codec, building Opus pipeline with timestamp handling...")
		args = append(args,
			"application/x-rtp,media=audio,encoding-name=OPUS,clock-rate=48000,payload=111", "!",
			"rtpjitterbuffer",
			"latency=0",
			"mode=none",
			"do-retransmission=false", "!",
			"rtpopusdepay", "!",
			"opusparse", "!",
			"webmmux", "!",
			"filesink", fmt.Sprintf("location=%s", outputFilePath),
		)
	} else if false {
		// Default to VP8 if codec is not detected
		r.logger.Info("Unknown or no codec detected, defaulting to VP8 pipeline with timestamp handling...")
		args = append(args,
			"application/x-rtp,media=video,encoding-name=VP8,clock-rate=90000", "!",
			"rtpjitterbuffer",
			"latency=0",
			"mode=none",
			"do-retransmission=false", "!",
			"rtpvp8depay", "!",
			"vp8parse", "!",
			"webmmux", "writing-app=GStreamer", "streamable=false", "min-index-interval=2000000000", "!",
			"filesink", fmt.Sprintf("location=%s", outputFilePath),
		)
	}

	r.logger.Info("GStreamer pipeline: %s", strings.Join(args, " ")) // Skip debug args for display

	r.gstreamerCmd = exec.Command("gst-launch-1.0", args...)
	// Redirect output for debugging
	r.gstreamerCmd.Stdout = os.Stdout
	r.gstreamerCmd.Stderr = os.Stderr

	// Start GStreamer process
	if err := r.gstreamerCmd.Start(); err != nil {
		return err
	}

	r.logger.Info("GStreamer pipeline started with PID: %d", r.gstreamerCmd.Process.Pid)

	// Monitor the process in a goroutine
	go func() {
		if err := r.gstreamerCmd.Wait(); err != nil {
			r.logger.Error("GStreamer process exited with error: %v", err)
		} else {
			r.logger.Info("GStreamer process exited normally")
		}
	}()

	return nil
}

// parseRtpCapsFromSDP extracts basic RTP caps from an SDP for use with application/x-rtp caps
// Prioritizes video codecs (H264/VP9/VP8/AV1) over audio (OPUS) and parses payload/clock-rate
func parseRtpCapsFromSDP(sdp string) (media string, encodingName string, payload int, clockRate int) {
	upper := strings.ToUpper(sdp)

	// Defaults
	media = "video"
	encodingName = "VP9"
	payload = 96
	clockRate = 90000

	// Select target encoding with priority: H264 > VP9 > VP8 > AV1 > OPUS (audio)
	if strings.Contains(upper, "H264") || strings.Contains(upper, "H.264") {
		encodingName = "H264"
		media = "video"
		clockRate = 90000
	} else if strings.Contains(upper, "VP9") {
		encodingName = "VP9"
		media = "video"
		clockRate = 90000
	} else if strings.Contains(upper, "VP8") {
		encodingName = "VP8"
		media = "video"
		clockRate = 90000
	} else if strings.Contains(upper, "AV1") {
		encodingName = "AV1"
		media = "video"
		clockRate = 90000
	} else if strings.Contains(upper, "OPUS") {
		encodingName = "OPUS"
		media = "audio"
		clockRate = 48000
	}

	// Parse matching a=rtpmap for the chosen encoding to refine payload and clock
	chosen := encodingName
	for _, line := range strings.Split(sdp, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(strings.ToLower(line), "a=rtpmap:") {
			continue
		}
		// Example: a=rtpmap:96 VP9/90000
		after := strings.TrimSpace(line[len("a=rtpmap:"):])
		fields := strings.Fields(after)
		if len(fields) < 2 {
			continue
		}
		ptStr := fields[0]
		codec := strings.ToUpper(fields[1])
		parts := strings.Split(codec, "/")
		name := parts[0]
		if name != chosen {
			continue
		}
		if v, err := strconv.Atoi(ptStr); err == nil {
			payload = v
		}
		if len(parts) >= 2 {
			if v, err := strconv.Atoi(parts[1]); err == nil {
				clockRate = v
			}
		}
		break
	}

	return
}

func (r *CursorGstreamerWebmRecorder) OnRTP(packet *rtp.Packet) error {
	// Marshal RTP packet
	buf, err := packet.Marshal()
	if err != nil {
		return err
	}

	return r.PushRtpBuf(buf)
}

func (r *CursorGstreamerWebmRecorder) PushRtpBuf(buf []byte) error {
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
			r.logger.Warn("Failed to write RTP length header: %v", err)
			return err
		}
		if _, err := r.rtpConn.Write(buf); err != nil {
			r.logger.Warn("Failed to write RTP packet: %v", err)
			return err
		}
	}
	return nil
}

func (r *CursorGstreamerWebmRecorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("Closing GStreamer WebM recorder...")

	r.logger.Info("Closing GStreamer WebM recorder2222...")

	// Cancel context to stop background goroutines
	if r.cancel != nil {
		r.cancel()
	}

	// Close TCP connection
	if r.rtpConn != nil {
		r.logger.Info("Closing TCP connection...")
		_ = r.rtpConn.Close()
		r.rtpConn = nil
		r.logger.Info("TCP connection closed")
	}

	// Gracefully stop GStreamer
	if r.gstreamerCmd != nil && r.gstreamerCmd.Process != nil {
		r.logger.Info("Stopping GStreamer process...")

		// Send EOS (End of Stream) signal to GStreamer
		// GStreamer handles SIGINT gracefully and will finish writing the file
		if err := r.gstreamerCmd.Process.Signal(os.Interrupt); err != nil {
			r.logger.Error("Failed to send SIGINT to GStreamer: %v", err)
			// If interrupt fails, force kill
			r.gstreamerCmd.Process.Kill()
		} else {
			r.logger.Info("Sent SIGINT to GStreamer, waiting for graceful exit...")

			// Wait for graceful exit with timeout
			done := make(chan error, 1)
			go func() {
				done <- r.gstreamerCmd.Wait()
			}()

			select {
			case <-time.After(15 * time.Second):
				r.logger.Info("GStreamer exit timeout, force killing...")
				// Timeout, force kill
				r.gstreamerCmd.Process.Kill()
				<-done // Wait for the kill to complete
			case err := <-done:
				if err != nil {
					r.logger.Info("GStreamer exited with error: %v", err)
				} else {
					r.logger.Info("GStreamer exited gracefully")
				}
			}
		}
	}

	// Clean up temporary SDP file
	if r.sdpFile != nil {
		os.Remove(r.sdpFile.Name())
		r.sdpFile = nil
	}

	// Post-process WebM to fix duration metadata if needed
	if r.tempOutputPath != "" && r.finalOutputPath != "" {
		r.logger.Info("Starting WebM duration post-processing...")
	}

	r.logger.Info("GStreamer WebM recorder closed")
	return nil
}

// GetOutputPath returns the output file path (for compatibility)
func (r *CursorGstreamerWebmRecorder) GetOutputPath() string {
	// Return final output path if post-processing is enabled, otherwise return original
	if r.finalOutputPath != "" {
		return r.finalOutputPath
	}
	return r.outputPath
}

// IsRecording returns true if the recorder is currently active
func (r *CursorGstreamerWebmRecorder) IsRecording() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.gstreamerCmd != nil && r.gstreamerCmd.Process != nil
}
