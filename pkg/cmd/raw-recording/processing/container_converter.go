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
	webrtc "github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/rtpdump"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"
)

const (
	audioMaxLate = 200  // 4sec
	videoMaxLate = 1000 // 4sec
)

type RTPDump2WebMConverter struct {
	logger        *ProcessingLogger
	reader        *rtpdump.Reader
	recorder      WebmRecorder
	sampleBuilder *samplebuilder.SampleBuilder

	lastPkt         *rtp.Packet
	lastPktDuration uint32
	dtxInserted     uint64

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

	for rtpdumpFile := range rtpdumpFiles {
		c := newRTPDump2WebMConverter(logger)
		if err := c.ConvertFile(rtpdumpFile, fixDtx); err != nil {
			c.logger.Errorf("Failed to convert %s: %v", rtpdumpFile, err)
			continue
		}
	}

	return nil
}

func (c *RTPDump2WebMConverter) ConvertFile(inputFile string, fixDtx bool) error {
	c.logger.Debugf("Converting %s", inputFile)

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

	sdpContent, _ := readSDP(strings.Replace(inputFile, suffixRtpDump, suffixSdp, 1))
	mType, _ := mimeType(sdpContent)
	_, suffix := outputFormatForMimeType(mType)

	switch mType {
	case webrtc.MimeTypeAV1:
		releasePacketHandler := samplebuilder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())
		c.sampleBuilder = samplebuilder.New(videoMaxLate, &codecs.AV1Depacketizer{}, 90000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(strings.Replace(inputFile, suffixRtpDump, suffix, 1), sdpContent, c.logger)
	case webrtc.MimeTypeVP9:
		releasePacketHandler := samplebuilder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())
		c.sampleBuilder = samplebuilder.New(videoMaxLate, &codecs.VP9Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(strings.Replace(inputFile, suffixRtpDump, suffix, 1), sdpContent, c.logger)
	case webrtc.MimeTypeH264:
		releasePacketHandler := samplebuilder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())
		c.sampleBuilder = samplebuilder.New(videoMaxLate, &codecs.H264Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(strings.Replace(inputFile, suffixRtpDump, suffix, 1), sdpContent, c.logger)
	case webrtc.MimeTypeVP8:
		releasePacketHandler := samplebuilder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())
		c.sampleBuilder = samplebuilder.New(videoMaxLate, &codecs.VP8Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(strings.Replace(inputFile, suffixRtpDump, suffix, 1), sdpContent, c.logger)
	case webrtc.MimeTypeOpus:
		releasePacketHandler := samplebuilder.WithPacketReleaseHandler(c.buildOpusReleasePacketHandler(fixDtx))
		c.sampleBuilder = samplebuilder.New(audioMaxLate, &codecs.OpusPacket{}, 48000, releasePacketHandler)
		c.recorder, err = NewGstreamerConverter(strings.Replace(inputFile, suffixRtpDump, suffix, 1), sdpContent, c.logger)
	default:
		return fmt.Errorf("unsupported codec type: %s", mType)
	}
	if err != nil {
		return fmt.Errorf("failed to create WebM recorder: %w", err)
	}
	defer func() {
		_ = c.recorder.Close()
	}()

	// Convert and feed RTP packets
	return c.feedPackets(mType, reader)
}

func (c *RTPDump2WebMConverter) feedPackets(mType string, reader *rtpdump.Reader) error {
	startTime := time.Now()

	i := uint64(0)
	for ; ; i++ {
		packet, err := reader.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return err
		} else if packet.IsRTCP {
			//			_ = c.recorder.PushRtcpBuf(packet.Payload)
			continue
		}

		// Unmarshal the RTP packet from the raw payload
		rtpPacket := &rtp.Packet{}
		if err := rtpPacket.Unmarshal(packet.Payload); err != nil {
			c.logger.Warnf("Failed to unmarshal RTP packet %d: %v", i, err)
			continue
		}

		// Push packet to samplebuilder for reordering
		c.sampleBuilder.Push(rtpPacket)

		// Log progress
		if i%10000 == 0 && i > 0 {
			c.logger.Debugf("Processed %d packets", i)
		}
	}

	if c.sampleBuilder != nil {
		c.sampleBuilder.Flush()
	}

	duration := time.Since(startTime).Round(time.Millisecond)

	c.logger.Infof("Finished feeding %d packets (%d dtxInserted, %d real) (frames: %d total, codec: %s) in %v ", i+c.dtxInserted, c.dtxInserted, i, c.totalFrames, mType, duration)

	return nil
}

func (c *RTPDump2WebMConverter) buildDefaultReleasePacketHandler() func(pkt *rtp.Packet) {
	return func(pkt *rtp.Packet) {
		if pkt.Marker {
			c.totalFrames++
		}

		if c.lastPkt != nil {
			if pkt.SequenceNumber-c.lastPkt.SequenceNumber > 1 {
				c.logger.Infof("Missing Packet Detected, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
			}
		}

		c.lastPkt = pkt

		if e := c.recorder.OnRTP(pkt); e != nil {
			c.logger.Warnf("Failed to record RTP packet SeqNum: %d RtpTs: %d: %v", pkt.SequenceNumber, pkt.Timestamp, e)
		}
	}
}

func (c *RTPDump2WebMConverter) buildOpusReleasePacketHandler(fixDtx bool) func(pkt *rtp.Packet) {
	return func(pkt *rtp.Packet) {
		pkt.SequenceNumber += uint16(c.dtxInserted)

		if c.lastPkt != nil {
			if pkt.SequenceNumber-c.lastPkt.SequenceNumber > 1 {
				c.logger.Infof("Missing Packet Detected, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
			}

			if fixDtx {
				tsDiff := c.timestampDiff(pkt.Timestamp, c.lastPkt.Timestamp)
				lastPktDuration := opusPacketDurationMs(c.lastPkt)
				rtpDuration := uint32(lastPktDuration * 48)

				if rtpDuration == 0 {
					rtpDuration = c.lastPktDuration
					c.logger.Infof("LastPacket with no duration, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
				} else {
					c.lastPktDuration = rtpDuration
				}

				if rtpDuration > 0 && tsDiff > rtpDuration {

					// Calculate how many packets we need to insert, taking care of packet losses
					var toAdd uint16
					if uint32(c.sequenceNumberDiff(pkt.SequenceNumber, c.lastPkt.SequenceNumber))*rtpDuration != tsDiff {
						toAdd = uint16(tsDiff/rtpDuration) - c.sequenceNumberDiff(pkt.SequenceNumber, c.lastPkt.SequenceNumber)
					}

					c.logger.Debugf("Gap detected, inserting %d packets tsDiff %d, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d",
						toAdd, tsDiff, c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)

					for i := 1; i <= int(toAdd); i++ {
						ins := c.lastPkt.Clone()
						ins.Payload = ins.Payload[:1] // Keeping only TOC byte
						ins.SequenceNumber += uint16(i)
						ins.Timestamp += uint32(i) * rtpDuration

						c.logger.Debugf("Writing dtxInserted Packet %v", ins)
						if e := c.recorder.OnRTP(ins); e != nil {
							c.logger.Warnf("Failed to record dtxInserted RTP packet SeqNum: %d RtpTs: %d: %v", ins.SequenceNumber, ins.Timestamp, e)
						}
					}

					c.dtxInserted += uint64(toAdd)
					pkt.SequenceNumber += toAdd
				}
			}
		}

		c.lastPkt = pkt

		c.logger.Debugf("Writing real Packet Last SeqNum: %d RtpTs: %d", pkt.SequenceNumber, pkt.Timestamp)
		if e := c.recorder.OnRTP(pkt); e != nil {
			c.logger.Warnf("Failed to record RTP packet SeqNum: %d RtpTs: %d: %v", pkt.SequenceNumber, pkt.Timestamp, e)
		}
	}
}

func getMaxFrameDimension(f1, f2 SegmentFrameDimension) SegmentFrameDimension {
	if f1.Width*f1.Height > f2.Width*f2.Height {
		return f1
	}
	return f2
}

func opusPacketDurationMs(pkt *rtp.Packet) int {
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
	frameDuration := float32(duration) / 10
	frameCount := opusFrameCount(c, payload)

	return int(frameDuration * float32(frameCount))
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

func (c *RTPDump2WebMConverter) timestampDiff(pts, fts uint32) uint32 {
	return pts - fts
}

func (c *RTPDump2WebMConverter) sequenceNumberDiff(psq, fsq uint16) uint16 {
	return psq - fsq
}
