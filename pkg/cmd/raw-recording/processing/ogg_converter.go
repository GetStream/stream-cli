package processing

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4/pkg/media/oggwriter"
)

const (
	// oggOpusChannels is the channel count written into OpusHead. WebRTC always
	// advertises Opus as "opus/48000/2" in the SDP rtpmap regardless of what the
	// encoder actually produces (the TOC byte carries the real configuration), so
	// this matches both the rtpmap and what rtpopusdepay+matroskamux emit on the
	// GStreamer path.
	oggOpusChannels = 2

	// oggWriteBufferSize batches Ogg pages before hitting the filesystem. The
	// writer emits one page per RTP packet, so an unbuffered sink costs one
	// write syscall per packet; batching cuts the write phase roughly in half.
	oggWriteBufferSize = 256 * 1024

	// oggOpusPreSkip is the OpusHead pre-skip oggwriter hardcodes, in 48kHz
	// samples. See writePriming.
	oggOpusPreSkip = 3840
)

// OggConverter writes Opus RTP packets straight into an Ogg Opus file, without
// handing them to an external process. It satisfies the same WebmRecorder
// contract as the GStreamer and FFmpeg converters.
//
// Unlike those, there is no socket and no second process between the feed loop
// and the file, so packets cannot be dropped in transit and the feed needs no
// pacing: a full buffer simply blocks the caller.
//
// Opus is the only codec this can serve — one RTP packet carries exactly one
// Opus packet, so no depacketiser is required. Fragmented video codecs still
// need GStreamer.
//
// The files it produces have no end-of-stream page: oggwriter only emits one
// when it owns the file descriptor, and taking that path would cost the write
// buffer and corrupt the tail whenever the final payload needs more than one
// lacing byte. These segments are intermediates consumed by a mix in the same
// run, and completeness is checked against the packet count, not the container.
//
// The timelines they declare are one packet duration early, and knowingly so:
// oggwriter stamps each page with the granule its packet *starts* at rather than
// ends at, and hardcodes the first page's granule to 1, which no timestamp we
// hand it can offset. A decoder therefore dates the first real sample at
// -(packetSamples-1) instead of 0, and a mix positions the segment that much
// early. It is uniform across segments, so relative sync between speakers is
// unaffected. Making it exact means owning the page layer (pre-skip 0, granule =
// packet end), which would also remove the need for writePriming entirely.
type OggConverter struct {
	logger     *ProcessingLogger
	outputPath string
	startAt    time.Time

	mu        sync.Mutex
	file      *os.File
	buf       *bufio.Writer
	writer    *oggwriter.OggWriter
	primed    bool
	closed    bool
	written   uint64
	skipped   uint64
	regressed uint64
	lastTS    uint32
}

// NewOggConverter creates an Ogg Opus writer for outputPath. sdpContent is used
// only to verify the track really is Opus at the expected clock rate; the Ogg
// headers are otherwise fully determined by the codec.
func NewOggConverter(outputPath, sdpContent string, logger *ProcessingLogger) (*OggConverter, error) {
	startAt := time.Now()

	_, encodingName, _, clockRate, err := parseRtpCapsFromSDP(sdpContent)
	if err != nil {
		return nil, err
	}
	if encodingName != "OPUS" {
		return nil, fmt.Errorf("ogg converter supports Opus only, got %s", encodingName)
	}
	rate, err := strconv.ParseUint(clockRate, 10, 32)
	if err != nil || rate == 0 {
		return nil, fmt.Errorf("invalid RTP clock-rate %q", clockRate)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create ogg output %s: %w", outputPath, err)
	}

	buf := bufio.NewWriterSize(file, oggWriteBufferSize)
	writer, err := oggwriter.NewWith(buf, uint32(rate), oggOpusChannels)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("failed to create ogg writer for %s: %w", outputPath, err)
	}

	logger.Infof("Ogg converter started for %s (%s/%s)", outputPath, encodingName, clockRate)

	return &OggConverter{
		logger:     logger,
		outputPath: outputPath,
		startAt:    startAt,
		file:       file,
		buf:        buf,
		writer:     writer,
	}, nil
}

// OnRTP appends one Opus packet to the Ogg stream. The page granule position is
// derived from the RTP timestamp, so DTX gaps and inserted packets carry through
// as real timing rather than being collapsed.
func (r *OggConverter) OnRTP(packet *rtp.Packet) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return fmt.Errorf("ogg converter for %s is closed", r.outputPath)
	}
	// oggwriter silently ignores empty payloads; count them here so the written
	// total stays a faithful record of what reached the file.
	if len(packet.Payload) == 0 {
		r.skipped++
		return nil
	}
	// A packet that does not advance the RTP timestamp cannot be given a page of
	// its own: oggwriter derives each page's granule position from the timestamp
	// delta, so the page would repeat the previous granule and the stream would
	// claim two packets at one instant. Senders do emit these — the first packet
	// of a talkspurt after DTX reuses the previous timestamp — and ffmpeg's
	// resampler answers the zero-length step by discarding a whole frame's worth
	// of samples, which produces a resample step with no output and aborts the
	// mix (assert in choose_input, up to at least ffmpeg 9).
	if r.written > 0 && int32(packet.Timestamp-r.lastTS) <= 0 {
		if r.regressed == 0 {
			r.logger.Warnf("Dropping RTP packet seq=%d with non-advancing timestamp %d (previous %d) from %s",
				packet.SequenceNumber, packet.Timestamp, r.lastTS, r.outputPath)
		}
		r.regressed++
		return nil
	}
	if !r.primed {
		r.primed = true
		r.writePriming(packet)
	}
	if err := r.writer.WriteRTP(packet); err != nil {
		return fmt.Errorf("failed to write RTP packet seq=%d into %s: %w", packet.SequenceNumber, r.outputPath, err)
	}
	r.written++
	r.lastTS = packet.Timestamp
	return nil
}

// writePriming emits exactly the pre-skip worth of throwaway audio ahead of the
// first real packet.
//
// oggwriter hardcodes the OpusHead pre-skip to 3840 samples and offers no way to
// set it, so a decoder always discards the first 80 ms of the stream. Our packets
// come off the wire mid-flight, so without priming that discard eats real audio
// and shifts the segment 80 ms early. Handing the decoder 80 ms to throw away —
// which is what pre-skip is for — keeps every captured sample.
// Two properties decide where the priming leaves the audio, and both were
// measured against 20/60/80ms fixtures rather than reasoned about:
//
//   - Total priming must equal the pre-skip exactly. Any surplus survives the
//     decoder's discard as leading silence and pushes the segment late by that
//     amount (60ms packets rounded up to 2x2880 measured +40ms).
//   - The packet immediately before the real audio must have the stream's own
//     duration. Priming a 60ms stream with 20ms packets covers the pre-skip
//     exactly yet still measured -40ms, because oggwriter derives each page's
//     granule from the RTP timestamp delta and the mismatched final step leaves
//     granule and decoded sample count disagreeing across the boundary.
//
// Both hold when the pre-skip splits as a frame-sized remainder followed by
// whole stream-shaped packets: 3840 = 2x480 + 1x2880 for a 60ms stream.
func (r *OggConverter) writePriming(first *rtp.Packet) {
	packetPayload, okPacket := opusSilencePacket(first.Payload)
	framePayload, frameSamples, okFrame := opusSingleFrameSilence(first.Payload)
	packetSamples := uint32(opusPacketDurationMsFactor10(first)*48) / 10

	if !okPacket || !okFrame || packetSamples == 0 || frameSamples == 0 {
		r.logger.Warnf("Cannot pad pre-skip for %s: no silence packet can be built from TOC 0x%02x, the first %d samples will be discarded by the decoder",
			r.outputPath, firstByte(first.Payload), oggOpusPreSkip)
		return
	}

	whole := oggOpusPreSkip / packetSamples
	remainder := uint32(oggOpusPreSkip) % packetSamples

	// A packet longer than the whole pre-skip cannot be preceded by anything, so
	// the two properties conflict; keeping the stream's shape wins and the
	// surplus is accepted (100ms costs 20ms of leading silence, 120ms costs 40ms).
	if whole == 0 {
		whole, remainder = 1, 0
	}
	if remainder%frameSamples != 0 {
		r.logger.Warnf("Cannot pad pre-skip for %s: %d-sample remainder is not a whole number of %d-sample frames, the first %d samples will be discarded by the decoder",
			r.outputPath, remainder, frameSamples, oggOpusPreSkip)
		return
	}

	// Oldest first, so the stream-shaped packets sit closest to the real audio.
	priming := *first
	ts := first.Timestamp - (whole*packetSamples + remainder)

	write := func(payload []byte, samples uint32) bool {
		priming.Payload = payload
		priming.Timestamp = ts
		if err := r.writer.WriteRTP(&priming); err != nil {
			r.logger.Warnf("Failed to write pre-skip priming packet for %s: %v", r.outputPath, err)
			return false
		}
		ts += samples
		return true
	}

	for i := uint32(0); i < remainder/frameSamples; i++ {
		if !write(framePayload, frameSamples) {
			return
		}
	}
	for i := uint32(0); i < whole; i++ {
		if !write(packetPayload, packetSamples) {
			return
		}
	}
}

// opusSingleFrameSilence returns a silent one-frame Opus packet in src's mode and
// bandwidth, plus that frame's length in 48kHz samples. The frame is the shortest
// its mode/bandwidth row offers — 10ms for SILK and Hybrid, 2.5ms for CELT — which
// is the finest padding granularity a stream can take without leaving its row.
//
// It is deliberately not a copy of src's own frame. 60ms is the one Opus frame
// duration that does not divide the pre-skip, so a 60ms stream is left 960 samples
// its own frame cannot subdivide and gets no priming at all. Every row's shortest
// frame does divide it (3840 = 8x480 = 32x120), and every longer frame in a row is
// that shortest one times a power of two, so the remainder always tiles.
func opusSingleFrameSilence(src []byte) ([]byte, uint32, bool) {
	if len(src) < 1 {
		return nil, 0, false
	}
	// Each row lists its durations ascending as the shortest x 2^n (RFC 6716
	// section 3.1), so the shortest is the row's base config. Only the two Hybrid
	// rows are two configs wide; the SILK and CELT rows are four.
	config := (src[0] >> 3) & 0x1F
	if config >= 12 && config < 16 {
		config &^= 0x01 // Hybrid SWB/FB: 10, 20 ms
	} else {
		config &^= 0x03 // SILK NB/MB/WB: 10..60 ms; CELT NB/WB/SWB/FB: 2.5..20 ms
	}

	// opusFrameDurationFactor10 returns tenths of a millisecond.
	tenths := opusFrameDurationFactor10(config)
	if tenths == 0 {
		return nil, 0, false
	}
	// Code 0 is exactly one frame, so the TOC alone is a valid zero-length frame.
	// The stereo bit carries over; src's frame-count code does not.
	return []byte{config<<3 | src[0]&0x04}, uint32(tenths) * 48 / 10, true
}

// PushRtpBuf unmarshals a marshalled RTP packet and appends it. OnRTP is the
// cheaper entry point; this exists to satisfy WebmRecorder.
func (r *OggConverter) PushRtpBuf(payload []byte) error {
	var packet rtp.Packet
	if err := packet.Unmarshal(payload); err != nil {
		return err
	}
	return r.OnRTP(&packet)
}

// Close finalises the Ogg stream and flushes it to disk. It is idempotent.
//
// Order matters: the buffer sits downstream of the Ogg writer, so the writer has
// to be finished before the buffer is flushed.
func (r *OggConverter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.closed {
		return nil
	}
	r.closed = true

	writerErr := r.writer.Close()
	flushErr := r.buf.Flush()
	closeErr := r.file.Close()

	r.logger.Infof("Ogg converter wrote %d packets (%d skipped, %d non-advancing) to %s in %s",
		r.written, r.skipped, r.regressed, r.outputPath, time.Since(r.startAt).Round(time.Millisecond))

	err := errors.Join(writerErr, flushErr, closeErr)
	if err != nil {
		return fmt.Errorf("failed to close ogg output %s: %w", r.outputPath, err)
	}
	return nil
}
