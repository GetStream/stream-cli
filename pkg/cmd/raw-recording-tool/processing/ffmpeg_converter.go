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

	"github.com/pion/rtcp"
	"github.com/pion/rtp"
)

type FfmpegConverter struct {
	logger     *ProcessingLogger
	outputPath string
	conn       *net.UDPConn
	ffmpegCmd  *exec.Cmd
	stdin      io.WriteCloser
	mu         sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	sdpFile    *os.File

	// Parsed from FFmpeg output: "Duration: N/A, start: <value>, bitrate: N/A"
	startOffsetMs  int64
	hasStartOffset bool
}

func NewFfmpegConverter(outputPath, sdpContent string, logger *ProcessingLogger) (*FfmpegConverter, error) {
	r := &FfmpegConverter{
		logger:     logger,
		outputPath: outputPath,
	}

	// Set up UDP connections
	port := rand.Intn(10000) + 10000
	if err := r.setupConnections(port); err != nil {
		return nil, err
	}

	// Start FFmpeg with codec detection
	if err := r.startFFmpeg(outputPath, sdpContent, port); err != nil {
		r.conn.Close()
		return nil, err
	}

	time.Sleep(2 * time.Second) // Wait for udp socket opened

	return r, nil
}

func (r *FfmpegConverter) setupConnections(port int) error {
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

	if e := r.conn.SetWriteBuffer(2 * 1024); e != nil {
		r.logger.Errorf("Failed to set UDP write buffer: %v", e)
	}
	if e := r.conn.SetReadBuffer(10 * 1024); e != nil {
		r.logger.Errorf("Failed to set UDP read buffer: %v", e)
	}

	return nil
}

func (r *FfmpegConverter) startFFmpeg(outputFilePath, sdpContent string, port int) error {

	// Write SDP to a temporary file
	sdpFile, err := os.CreateTemp("", "cursor_webm_*.sdp")
	if err != nil {
		return err
	}
	r.sdpFile = sdpFile

	updatedSdp := replaceSDP(sdpContent, port)
	r.logger.Infof("Using Sdp:\n%s\n", updatedSdp)

	if _, e := r.sdpFile.WriteString(updatedSdp); e != nil {
		r.sdpFile.Close()
		return e
	}
	r.sdpFile.Close()

	// Build FFmpeg command with optimized settings for single track recording
	args := r.generateArgs(sdpFile.Name(), outputFilePath)

	r.logger.Infof("FFMpeg pipeline: %s", strings.Join(args, " ")) // Skip debug args for display

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
	if e := r.ffmpegCmd.Start(); e != nil {
		return e
	}

	return nil
}

func (r *FfmpegConverter) generateArgs(sdp, output string) []string {
	// Build FFmpeg command with optimized settings for single track recording
	var args []string
	args = append(args, "-hide_banner")
	args = append(args, "-threads", "1")
	args = append(args, "-protocol_whitelist", "file,udp,rtp")
	args = append(args, "-buffer_size", "10000000")
	args = append(args, "-max_delay", "1000000")
	args = append(args, "-reorder_queue_size", "0")
	args = append(args, "-i", sdp)
	args = append(args, "-c", "copy")
	args = append(args, "-y", output)
	return args
}

// scanFFmpegOutput reads lines from FFmpeg output, mirrors to console, and extracts start offset.
func (r *FfmpegConverter) scanFFmpegOutput(reader io.Reader, isStderr bool) {
	scanner := bufio.NewScanner(reader)
	re := regexp.MustCompile(`\bstart:\s*([0-9]+(?:\.[0-9]+)?)`)
	for scanner.Scan() {
		line := scanner.Text()
		// Mirror output
		if isStderr {
			_, _ = fmt.Fprintln(os.Stderr, line)
		} else {
			_, _ = fmt.Fprintln(os.Stdout, line)
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
					r.logger.Infof("Detected FFmpeg start offset: %d ms", r.startOffsetMs)
				}
				r.mu.Unlock()
			}
		}
	}
	_ = scanner.Err()
}

// StartOffset returns the parsed FFmpeg start offset in seconds and whether it was found.
func (r *FfmpegConverter) StartOffset() (int64, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.startOffsetMs, r.hasStartOffset
}

func (r *FfmpegConverter) OnRTP(packet *rtp.Packet) error {
	// Marshal RTP packet
	buf, err := packet.Marshal()
	if err != nil {
		return err
	}

	return r.PushRtpBuf(buf)
}

func (r *FfmpegConverter) PushRtpBuf(buf []byte) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Send RTP packet over UDP
	if r.conn != nil {
		_, err := r.conn.Write(buf)
		if err != nil {
			r.logger.Warnf("Failed to write RTP packet: %v", err)
		}
		return err
	}
	return nil
}

func (r *FfmpegConverter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Cancel context to stop background goroutines
	if r.cancel != nil {
		r.cancel()
	}

	r.logger.Infof("Closing UPD connection and wait for FFMpeg termination...")

	// Close UDP connection by sending arbitrary RtcpBye (Ffmpeg is now able to end correctly)
	if r.conn != nil {
		buf, _ := rtcp.Goodbye{
			Sources: []uint32{1}, // fixed ssrc is ok
			Reason:  "bye",
		}.Marshal()
		_, _ = r.conn.Write(buf)
		_ = r.conn.Close()
		r.conn = nil
	}

	// Gracefully wait for FFmpeg termination
	if r.ffmpegCmd != nil && r.ffmpegCmd.Process != nil {
		// Wait for graceful exit with timeout
		done := make(chan error, 1)
		go func() {
			done <- r.ffmpegCmd.Wait()
		}()

		select {
		case <-time.After(5 * time.Second):
			r.logger.Warnf("FFMpeg Process termination timeout...")

			// Timeout, force kill
			if e := r.ffmpegCmd.Process.Kill(); e != nil {
				r.logger.Errorf("FFMpeg Process errored while killing: %v", e)
			}
		case <-done:
			r.logger.Infof("FFMpeg Process exited succesfully...")
		}
	}

	// Clean up temporary SDP file
	if r.sdpFile != nil {
		_ = os.Remove(r.sdpFile.Name())
		r.sdpFile = nil
	}

	return nil
}
