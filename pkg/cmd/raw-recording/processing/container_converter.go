package processing

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4/pkg/media/rtpdump"

	"github.com/GetStream/stream-cli/pkg/cmd/raw-recording/processing/rtpreorder"
)

const (
	audioMaxLate = 1000 // 20sec
	videoMaxLate = 5000 // 20sec
)

type RTPDump2WebMConverter struct {
	logger    *ProcessingLogger
	reader    *rtpdump.Reader
	recorder  WebmRecorder
	reorderer *rtpreorder.Reorderer
	recordErr error

	lastPkt         *rtp.Packet
	lastPktDuration uint32
	dtxInserted     uint64
	missingPkt      uint64

	totalFrames int
}

type WebmRecorder interface {
	OnRTP(pkt *rtp.Packet) error
	PushRtpBuf(payload []byte) error
	Close() error
}

func newRTPDump2WebMConverter(logger *ProcessingLogger) *RTPDump2WebMConverter {
	return &RTPDump2WebMConverter{
		logger: logger,
	}
}

func ConvertDirectory(directory string, accept func(path string, info os.FileInfo) (*SegmentInfo, bool), fixDtx bool, logger *ProcessingLogger) error {
	rtpdumpFiles := make(map[string]*SegmentInfo)

	// Walk through directory to find .rtpdump files
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), suffixRtpDump) {
			segment, accepted := accept(path, info)
			if accepted {
				rtpdumpFiles[path] = segment
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	for rtpdumpFile, seg := range rtpdumpFiles {
		c := newRTPDump2WebMConverter(logger)
		if err := c.ConvertFile(rtpdumpFile, seg, fixDtx); err != nil {
			c.logger.Errorf("Failed to convert %s: %v", rtpdumpFile, err)
			continue
		}
	}

	return nil
}

func (c *RTPDump2WebMConverter) ConvertFile(inputFile string, segment *SegmentInfo, fixDtx bool) (err error) {
	conversionStart := time.Now()
	defer func() {
		if err != nil {
			c.logger.Infof("Conversion failed after %s for %s: %v", time.Since(conversionStart).Round(time.Millisecond), inputFile, err)
			return
		}
		c.logger.Infof("Conversion completed in %s for %s", time.Since(conversionStart).Round(time.Millisecond), inputFile)
	}()

	c.logger.Infof("Converting %s - FixDtx %t", inputFile, fixDtx)

	// Parse the RTP dump file
	// Open the file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open rtpdump file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	// Create standardized reader
	reader, _, _ := rtpdump.NewReader(file)
	c.reader = reader

	sdpContent, err := readSDP(strings.Replace(inputFile, suffixRtpDump, suffixSdp, 1))
	if err != nil {
		return err
	}
	mType, err := mimeType(sdpContent)
	if err != nil {
		return err
	}

	switch mType {
	case webrtc.MimeTypeAV1:
		releasePacketHandler := rtpreorder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())
		c.reorderer = rtpreorder.New(videoMaxLate, &codecs.AV1Depacketizer{}, 90000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(segment.ContainerPath, sdpContent, c.logger)
	case webrtc.MimeTypeVP9:
		releasePacketHandler := rtpreorder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())
		c.reorderer = rtpreorder.New(videoMaxLate, &codecs.VP9Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(segment.ContainerPath, sdpContent, c.logger)
	case webrtc.MimeTypeH264:
		releasePacketHandler := rtpreorder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())
		c.reorderer = rtpreorder.New(videoMaxLate, &codecs.H264Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(segment.ContainerPath, sdpContent, c.logger)
	case webrtc.MimeTypeVP8:
		releasePacketHandler := rtpreorder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())
		c.reorderer = rtpreorder.New(videoMaxLate, &codecs.VP8Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(segment.ContainerPath, sdpContent, c.logger)
	case webrtc.MimeTypeOpus:
		releasePacketHandler := rtpreorder.WithPacketReleaseHandler(c.buildOpusReleasePacketHandler(fixDtx))
		c.reorderer = rtpreorder.New(audioMaxLate, &codecs.OpusPacket{}, 48000, releasePacketHandler)
		// One RTP packet is one Opus packet, so an Ogg target needs no
		// depacketiser and can be written in-process.
		if segment.ContainerExt == oggExtension {
			c.recorder, err = NewOggConverter(segment.ContainerPath, sdpContent, c.logger)
		} else {
			c.recorder, err = NewGstreamerConverter(segment.ContainerPath, sdpContent, c.logger)
		}
	default:
		return fmt.Errorf("unsupported codec type: %s", mType)
	}
	if err != nil {
		return fmt.Errorf("failed to create WebM recorder: %w", err)
	}
	defer func() {
		if e := recover(); e != nil {
			_ = c.recorder.Close()
			panic(e)
		}
	}()

	// Convert and feed RTP packets
	feedErr := c.feedPackets(mType, reader)
	closeErr := c.recorder.Close()
	if feedErr != nil && closeErr != nil {
		c.logger.Warnf("Recorder close failed after feed error for %s: %v", inputFile, closeErr)
	}
	if feedErr != nil {
		return feedErr
	}
	if closeErr != nil {
		return closeErr
	}
	return nil
}

func (c *RTPDump2WebMConverter) feedPackets(mType string, reader *rtpdump.Reader) error {
	startTime := time.Now()

	var i, pushed, rtcp uint64
	for ; ; i++ {
		packet, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		}

		if packet.IsRTCP {
			rtcp++
			//			_ = c.recorder.PushRtcpBuf(packet.Payload)
			continue
		}

		// Unmarshal the RTP packet from the raw payload
		rtpPacket := &rtp.Packet{}
		if err := rtpPacket.Unmarshal(packet.Payload); err != nil {
			c.logger.Warnf("Failed to unmarshal RTP packet %d: %v", i, err)
			continue
		}

		// Push packet to the reorderer before forwarding to the recorder.
		pushed++
		c.reorderer.Push(rtpPacket)

		// Log progress
		if (i+c.dtxInserted)%50000 == 0 && i > 0 {
			c.logger.Infof("Processed %d packets", i)
		}
	}

	if c.reorderer != nil {
		c.reorderer.Flush()
	}

	duration := time.Since(startTime).Round(time.Millisecond)

	c.logger.Infof("Finished feeding %d RTP packets (%d dtxInserted, %d recorded, %d missing, %d rtcp) (frames: %d total, codec: %s) in %v", pushed+c.dtxInserted, c.dtxInserted, i, c.missingPkt, rtcp, c.totalFrames, mType, duration)

	if c.recordErr != nil {
		return c.recordErr
	}
	return nil
}

func (c *RTPDump2WebMConverter) buildDefaultReleasePacketHandler() func(pkt *rtp.Packet) {
	return func(pkt *rtp.Packet) {
		if pkt.Marker {
			c.totalFrames++
		}

		if c.lastPkt != nil {
			diff := pkt.SequenceNumber - c.lastPkt.SequenceNumber
			if diff > 1 {
				c.missingPkt += uint64(diff - 1)
				c.logger.Debugf("Missing %d Packet Detected, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", pkt.SequenceNumber-c.lastPkt.SequenceNumber, c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
			}
		}

		c.lastPkt = pkt

		c.recordPacket(pkt, "RTP")
	}
}

func (c *RTPDump2WebMConverter) buildOpusReleasePacketHandler(fixDtx bool) func(pkt *rtp.Packet) {
	return func(pkt *rtp.Packet) {
		pkt.SequenceNumber += uint16(c.dtxInserted)

		if c.lastPkt != nil {
			diff := pkt.SequenceNumber - c.lastPkt.SequenceNumber
			if diff > 1 {
				c.missingPkt += uint64(diff - 1)
				c.logger.Debugf("Missing Packet Detected, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
			}

			if fixDtx {
				tsDiff := c.timestampDiff(pkt.Timestamp, c.lastPkt.Timestamp)
				rtpDuration := uint32(opusPacketDurationMsFactor10(c.lastPkt)*48) / 10

				if rtpDuration == 0 {
					rtpDuration = c.lastPktDuration
					c.logger.Infof("LastPacket with no duration, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
				} else {
					c.lastPktDuration = rtpDuration
				}

				if rtpDuration > 0 && tsDiff > rtpDuration {

					// Calculate how many packets we need to insert, taking care of packet losses
					toAdd := opusDtxPacketsToInsert(
						tsDiff,
						rtpDuration,
						c.sequenceNumberDiff(pkt.SequenceNumber, c.lastPkt.SequenceNumber),
					)

					c.logger.Debugf("Gap detected, inserting %d packets tsDiff %d, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d",
						toAdd, tsDiff, c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)

					silence, ok := opusSilencePacket(c.lastPkt.Payload)
					if !ok {
						// Emitting a malformed packet would corrupt the whole
						// container, which is far worse than leaving the hole:
						// ffmpeg abandons demuxing the stream entirely. Leave it
						// for the mixer's gap-fill filter and carry on.
						c.logger.Warnf("Cannot build a silence packet for a %d-byte Opus payload (TOC 0x%02x): leaving a %dms gap unfilled",
							len(c.lastPkt.Payload), firstByte(c.lastPkt.Payload), tsDiff/48)
						toAdd = 0
					}

					for i := 1; i <= int(toAdd); i++ {
						ins := c.lastPkt.Clone()
						ins.Payload = silence
						ins.SequenceNumber += uint16(i)
						ins.Timestamp += uint32(i) * rtpDuration

						c.logger.Debugf("Writing dtxInserted Packet %v", ins)
						c.recordPacket(ins, "dtxInserted RTP")
					}

					c.dtxInserted += uint64(toAdd)
					pkt.SequenceNumber += uint16(toAdd)
				}
			}
		}

		c.lastPkt = pkt

		c.recordPacket(pkt, "RTP")
	}
}

func (c *RTPDump2WebMConverter) recordPacket(pkt *rtp.Packet, label string) {
	if err := c.recorder.OnRTP(pkt); err != nil {
		c.logger.Warnf("Failed to record %s packet SeqNum: %d RtpTs: %d: %v", label, pkt.SequenceNumber, pkt.Timestamp, err)
		if c.recordErr == nil {
			c.recordErr = err
		}
	}
}

func getMaxFrameDimension(f1, f2 SegmentFrameDimension) SegmentFrameDimension {
	if f1.Width*f1.Height > f2.Width*f2.Height {
		return f1
	}
	return f2
}

func opusPacketDurationMsFactor10(pkt *rtp.Packet) int {
	// 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// | config  |s|1|1|0|p|     M     |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	payload := pkt.Payload
	if len(payload) < 1 {
		return 0
	}

	toc := payload[0]
	config := (toc >> 3) & 0x1F
	c := toc & 0x03

	// Calculate frame duration according to OPUS RFC 6716 table (use x10 factor)
	// Frame duration is determined by the config value
	duration := opusFrameDurationFactor10(config)
	frameCount := opusFrameCount(c, payload)

	return duration * frameCount
}

func opusFrameDurationFactor10(config byte) int {
	switch {
	case config < 3:
		// SILK-only NB: 10, 20, 40 ms
		return 100 * (1 << (config & 0x03))
	case config == 3:
		// SILK-only NB: 60 ms
		return 600
	case config < 7:
		// SILK-only MB: 10, 20, 40 ms
		return 100 * (1 << (config & 0x03))
	case config == 7:
		// SILK-only MB: 60 ms
		return 600
	case config < 11:
		// SILK-only WB: 10, 20, 40 ms
		return 100 * (1 << (config & 0x03))
	case config == 11:
		// SILK-only WB: 60 ms
		return 600
	case config <= 13:
		// Hybrid SWB: 10, 20 ms
		return 100 * (1 << (config & 0x01))
	case config <= 15:
		// Hybrid FB: 10, 20 ms
		return 100 * (1 << (config & 0x01))
	case config <= 19:
		// CELT-only NB: 2.5, 5, 10, 20 ms
		return 25 * (1 << (config & 0x03)) // 2.5ms * 10 for integer math
	case config <= 23:
		// CELT-only WB: 2.5, 5, 10, 20 ms
		return 25 * (1 << (config & 0x03)) // 2.5ms * 10 for integer math
	case config <= 27:
		// CELT-only SWB: 2.5, 5, 10, 20 ms
		return 25 * (1 << (config & 0x03)) // 2.5ms * 10 for integer math
	case config <= 31:
		// CELT-only FB: 2.5, 5, 10, 20 ms
		return 25 * (1 << (config & 0x03)) // 2.5ms * 10 for integer math
	default:
		// MUST NOT HAPPEN
		return 0
	}
}

func opusFrameCount(c byte, payload []byte) int {
	switch c {
	case 0:
		return 1
	case 1, 2:
		return 2
	case 3:
		if len(payload) > 1 {
			return int(payload[1] & 0x3F)
		}
	}
	return 0
}

func (c *RTPDump2WebMConverter) offset(pts, fts uint32) int64 {
	return int64(c.timestampDiff(pts, fts) / 90)
}

func (c *RTPDump2WebMConverter) timestampDiff(pts, fts uint32) uint32 {
	return pts - fts
}

func (c *RTPDump2WebMConverter) sequenceNumberDiff(psq, fsq uint16) uint16 {
	return psq - fsq
}

func firstByte(b []byte) byte {
	if len(b) == 0 {
		return 0
	}
	return b[0]
}

// opusSilencePacket builds the smallest structurally valid Opus packet carrying
// the same frame layout — and therefore the same duration — as src, but no audio
// data. It is the packet used to fill a timestamp hole.
//
// The frame count is encoded differently per TOC code, and only codes 0 and 1
// carry it in the TOC byte alone. Truncating every packet to its TOC (what this
// used to do) therefore produces a malformed packet for codes 2 and 3, and a
// malformed packet is not a local defect: ffmpeg reports "Invalid data found
// when processing input" and abandons the stream, so one bad filler destroys the
// entire recording.
//
// Code 3 is not exotic. Any sender packing multiple frames per packet uses it —
// a 60ms ptime arrives as one code-3 packet holding 3x20ms frames, with the
// frame count in payload byte 1.
//
// See RFC 6716 section 3.1 (TOC byte) and 3.2 (frame packing).
func opusSilencePacket(src []byte) ([]byte, bool) {
	if len(src) < 1 {
		return nil, false
	}
	toc := src[0]

	switch toc & 0x03 {
	case 0:
		// One frame; the TOC alone describes it, so a bare TOC is one
		// zero-length frame.
		return []byte{toc}, true
	case 1:
		// Two frames of equal size. The remaining length (zero) splits evenly.
		return []byte{toc}, true
	case 2:
		// Two frames of different sizes: the first frame's length is mandatory.
		// Zero leaves both frames empty.
		return []byte{toc, 0}, true
	case 3:
		// The frame count lives in byte 1 and must survive. Clearing the VBR and
		// padding bits leaves M equal-length frames sharing zero bytes of data.
		if len(src) < 2 {
			return nil, false
		}
		frames := src[1] & 0x3F
		if frames == 0 {
			return nil, false
		}
		return []byte{toc, frames}, true
	}

	return nil, false
}

func opusDtxPacketsToInsert(tsDiff, rtpDuration uint32, seqDiff uint16) uint32 {
	if rtpDuration == 0 {
		return 0
	}

	expectedPacketDiff := tsDiff / rtpDuration
	observedPacketDiff := uint32(seqDiff)
	if expectedPacketDiff <= observedPacketDiff {
		return 0
	}
	return expectedPacketDiff - observedPacketDiff
}
