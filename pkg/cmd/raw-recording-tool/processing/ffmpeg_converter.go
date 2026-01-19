package processing

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/GetStream/getstream-go/v3"
	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

type CursorWebmRecorder struct {
	logger     *getstream.DefaultLogger
	outputPath string
	conn       *net.UDPConn
	ffmpegCmd  *exec.Cmd
	stdin      io.WriteCloser
	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc

	// Parsed from FFmpeg output: "Duration: N/A, start: <value>, bitrate: N/A"
	startOffsetMs  int64
	hasStartOffset bool
}

func NewCursorWebmRecorder(outputPath, sdpContent string, logger *getstream.DefaultLogger) (*CursorWebmRecorder, error) {
	ctx, cancel := context.WithCancel(context.Background())

	r := &CursorWebmRecorder{
		logger:     logger,
		outputPath: outputPath,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Set up UDP connections
	port := rand.Intn(10000) + 10000
	if err := r.setupConnections(port); err != nil {
		cancel()
		return nil, err
	}

	// Start FFmpeg with codec detection
	if err := r.startFFmpeg(outputPath, sdpContent, port); err != nil {
		cancel()
		return nil, err
	}

	return r, nil
}

func (r *CursorWebmRecorder) setupConnections(port int) error {
	// Setup connection
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	r.conn = conn

	if e := r.conn.SetWriteBuffer(2048); e != nil {
		r.logger.Error("Failed to set UDP write buffer: %v", e)
	}
	if e := r.conn.SetReadBuffer(2048); e != nil {
		r.logger.Error("Failed to set UDP read buffer: %v", e)
	}

	return nil
}

func (r *CursorWebmRecorder) startFFmpeg(outputFilePath, sdpContent string, port int) error {

	// Write SDP to a temporary file
	sdpFile, err := os.CreateTemp("", "cursor_webm_*.sdp")
	if err != nil {
		return err
	}

	updatedSdp := replaceSDP(sdpContent, port)
	r.logger.Info("Using Sdp:\n%s\n", updatedSdp)

	if _, err := sdpFile.WriteString(updatedSdp); err != nil {
		sdpFile.Close()
		return err
	}
	sdpFile.Close()

	// Build FFmpeg command with optimized settings for single track recording
	args := []string{
		"-threads", "1",
		//		"-loglevel", "debug",
		"-protocol_whitelist", "file,udp,rtp",
		"-buffer_size", "425984",
		"-max_delay", "150000",
		"-reorder_queue_size", "0",
		"-i", sdpFile.Name(),
	}

	//switch strings.ToLower(mimeType) {
	//case "audio/opus":
	//	// For other codecs, use direct copy
	args = append(args, "-c", "copy")
	//default:
	//	// For other codecs, use direct copy
	//	args = append(args, "-c", "copy")
	//}
	//if isVP9 {
	//	// For VP9, avoid direct copy and use re-encoding with error resilience
	//	// This works around FFmpeg's experimental VP9 RTP support issues
	//	r.logger.Info("Detected VP9 codec, applying workarounds...")
	//	args = append(args,
	//		"-c:v", "libvpx-vp9",
	//		//			"-error_resilience", "aggressive",
	//		"-err_detect", "ignore_err",
	//		"-fflags", "+genpts+igndts",
	//		"-avoid_negative_ts", "make_zero",
	//		// VP9-specific quality settings to handle corrupted frames
	//		"-crf", "30",
	//		"-row-mt", "1",
	//		"-frame-parallel", "1",
	//	)
	//} else if strings.Contains(strings.ToUpper(sdpContent), "AV1") {
	//	args = append(args,
	//		"-c:v", "libaom-av1",
	//		"-cpu-used", "8",
	//		"-usage", "realtime",
	//	)
	//} else if strings.Contains(strings.ToUpper(sdpContent), "OPUS") {
	//	args = append(args, "-fflags", "+genpts", "-use_wallclock_as_timestamps", "0", "-c:a", "copy")
	//} else {
	//	// For other codecs, use direct copy
	//	args = append(args, "-c", "copy")
	//}

	args = append(args,
		"-y",
		outputFilePath,
	)

	r.logger.Info("FFMpeg pipeline: %s", strings.Join(args, " ")) // Skip debug args for display

	r.ffmpegCmd = exec.Command("ffmpeg", args...)

	// Capture stdout/stderr to parse FFmpeg logs while mirroring to console
	stdoutPipe, err := r.ffmpegCmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := r.ffmpegCmd.StderrPipe()
	if err != nil {
		return err
	}

	// Create stdin pipe to send commands to FFmpeg
	//var err error
	r.stdin, err = r.ffmpegCmd.StdinPipe()
	if err != nil {
		fmt.Println("Error creating stdin pipe:", err)
	}

	// Begin scanning output streams after process has started
	go r.scanFFmpegOutput(stdoutPipe, false)
	go r.scanFFmpegOutput(stderrPipe, true)

	// Start FFmpeg process
	if err := r.ffmpegCmd.Start(); err != nil {
		return err
	}

	return nil
}

// scanFFmpegOutput reads lines from FFmpeg output, mirrors to console, and extracts start offset.
func (r *CursorWebmRecorder) scanFFmpegOutput(reader io.Reader, isStderr bool) {
	scanner := bufio.NewScanner(reader)
	re := regexp.MustCompile(`\bstart:\s*([0-9]+(?:\.[0-9]+)?)`)
	for scanner.Scan() {
		line := scanner.Text()
		// Mirror output
		if isStderr {
			fmt.Fprintln(os.Stderr, line)
		} else {
			fmt.Fprintln(os.Stdout, line)
		}

		// Try to extract the start value from those lines  "Duration: N/A, start: 0.000000, bitrate: N/A"
		if !strings.Contains(line, "Duration") || !strings.Contains(line, "bitrate") {
			continue
		} else if matches := re.FindStringSubmatch(line); len(matches) == 2 {
			if v, parseErr := strconv.ParseFloat(matches[1], 64); parseErr == nil {
				// Save only once
				r.mu.Lock()
				if !r.hasStartOffset {
					r.startOffsetMs = int64(v * 1000)
					r.hasStartOffset = true
					r.logger.Info("Detected FFmpeg start offset: %.6f seconds", v)
				}
				r.mu.Unlock()
			}
		}
	}
	_ = scanner.Err()
}

// StartOffset returns the parsed FFmpeg start offset in seconds and whether it was found.
func (r *CursorWebmRecorder) StartOffset() (int64, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.startOffsetMs, r.hasStartOffset
}

func (r *CursorWebmRecorder) OnRTP(packet *rtp.Packet) error {
	// Marshal RTP packet
	buf, err := packet.Marshal()
	if err != nil {
		return err
	}

	return r.PushRtpBuf(buf)
}

func (r *CursorWebmRecorder) PushRtpBuf(buf []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Send RTP packet over UDP
	if r.conn != nil {
		r.conn.SetWriteDeadline(time.Now().Add(1000 * time.Microsecond))
		_, err := r.conn.Write(buf)
		if err != nil {
			//	return err)
			//}
			r.logger.Info("Wrote packet to %s - %v", r.conn.LocalAddr().String(), err)
		}
	}
	return nil
}

func (r *CursorWebmRecorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Cancel context to stop background goroutines
	if r.cancel != nil {
		r.cancel()
	}

	r.logger.Info("Closing UPD connection...")

	// Close UDP connection by sending arbitrary RtcpBye (Ffmpeg is no able to end correctly)
	if r.conn != nil {
		buf, _ := rtcp.Goodbye{
			Sources: []uint32{1}, // fixed ssrc is ok
			Reason:  "bye",
		}.Marshal()
		_, _ = r.conn.Write(buf)
		_ = r.conn.Close()
		r.conn = nil
	}

	r.logger.Info("UDP Connection closed...")

	time.Sleep(5 * time.Second)

	r.logger.Info("After sleep...")

	// Gracefully stop FFmpeg
	if r.ffmpegCmd != nil && r.ffmpegCmd.Process != nil {

		// ✅ Gracefully stop FFmpeg by sending 'q' to stdin
		//fmt.Println("Sending 'q' to FFmpeg...")
		//_, _ = r.stdin.Write([]byte("q\n"))
		//r.stdin.Close()

		// Send interrupt signal to FFmpeg process
		r.logger.Info("Sending SIGTERM...")

		//if err := r.ffmpegCmd.Process.Signal(os.Interrupt); err != nil {
		//	// If interrupt fails, force kill
		//	r.ffmpegCmd.Process.Kill()
		//} else {

		r.logger.Info("Waiting for SIGTERM...")

		// Wait for graceful exit with timeout
		done := make(chan error, 1)
		go func() {
			done <- r.ffmpegCmd.Wait()
		}()

		select {
		case <-time.After(10 * time.Second):
			r.logger.Info("Wait timetout for SIGTERM...")

			// Timeout, force kill
			r.ffmpegCmd.Process.Kill()
		case <-done:
			r.logger.Info("Process exited succesfully SIGTERM...")
			// Process exited gracefully
		}
	}

	return nil
}
