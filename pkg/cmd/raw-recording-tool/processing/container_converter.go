package processing

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GetStream/getstream-go/v3"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/rtpdump"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"
)

const audioMaxLate = 200  // 4sec
const videoMaxLate = 1000 // 4sec

type RTPDump2WebMConverter struct {
	logger        *getstream.DefaultLogger
	reader        *rtpdump.Reader
	recorder      WebmRecorder
	sampleBuilder *samplebuilder.SampleBuilder

	lastPkt         *rtp.Packet
	lastPktDuration uint32
	inserted        uint16
}

type WebmRecorder interface {
	OnRTP(pkt *rtp.Packet) error
	PushRtpBuf(payload []byte) error
	Close() error
}

func newRTPDump2WebMConverter(logger *getstream.DefaultLogger) *RTPDump2WebMConverter {
	return &RTPDump2WebMConverter{
		logger: logger,
	}
}

func ConvertDirectory(directory string, accept func(path string, info os.FileInfo) (*SegmentInfo, bool), fixDtx bool, logger *getstream.DefaultLogger) error {
	rtpdumpFiles := make(map[string]*SegmentInfo)

	// Walk through directory to find .rtpdump files
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), SuffixRtpDump) {
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

	for rtpdumpFile, segment := range rtpdumpFiles {
		c := newRTPDump2WebMConverter(logger)
		if err := c.ConvertFile(rtpdumpFile, fixDtx); err != nil {
			c.logger.Error("Failed to convert %s: %v", rtpdumpFile, err)
			continue
		}

		switch c.recorder.(type) {
		case *CursorWebmRecorder:
			offset, exists := c.recorder.(*CursorWebmRecorder).StartOffset()
			if exists {
				segment.FFMpegOffset = offset
			}
		}
	}

	return nil
}

func (c *RTPDump2WebMConverter) ConvertFile(inputFile string, fixDtx bool) error {
	c.logger.Info("Converting %s", inputFile)

	// Parse the RTP dump file
	// Open the file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open rtpdump file: %w", err)
	}
	defer file.Close()

	// Create standardized reader
	reader, _, _ := rtpdump.NewReader(file)
	c.reader = reader

	sdpContent, _ := readSDP(strings.Replace(inputFile, SuffixRtpDump, SuffixSdp, 1))
	mType, _ := mimeType(sdpContent)

	releasePacketHandler := samplebuilder.WithPacketReleaseHandler(c.buildDefaultReleasePacketHandler())

	switch mType {
	case webrtc.MimeTypeAV1:
		c.sampleBuilder = samplebuilder.New(videoMaxLate, &codecs.AV1Depacketizer{}, 90000, releasePacketHandler)
		c.recorder, err = NewCursorWebmRecorder(strings.Replace(inputFile, SuffixRtpDump, SuffixWebm, 1), sdpContent, c.logger)
	case webrtc.MimeTypeVP9:
		c.sampleBuilder = samplebuilder.New(videoMaxLate, &codecs.VP9Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewCursorGstreamerWebmRecorder(strings.Replace(inputFile, SuffixRtpDump, SuffixWebm, 1), sdpContent, c.logger)
	case webrtc.MimeTypeH264:
		c.sampleBuilder = samplebuilder.New(videoMaxLate, &codecs.H264Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewCursorWebmRecorder(strings.Replace(inputFile, SuffixRtpDump, SuffixMp4, 1), sdpContent, c.logger)
	case webrtc.MimeTypeVP8:
		c.sampleBuilder = samplebuilder.New(videoMaxLate, &codecs.VP8Packet{}, 90000, releasePacketHandler)
		c.recorder, err = NewCursorWebmRecorder(strings.Replace(inputFile, SuffixRtpDump, SuffixWebm, 1), sdpContent, c.logger)
	case webrtc.MimeTypeOpus:
		if fixDtx {
			releasePacketHandler = samplebuilder.WithPacketReleaseHandler(c.buildOpusReleasePacketHandler())
		}
		c.sampleBuilder = samplebuilder.New(audioMaxLate, &codecs.OpusPacket{}, 48000, releasePacketHandler)
		c.recorder, err = NewCursorWebmRecorder(strings.Replace(inputFile, SuffixRtpDump, SuffixWebm, 1), sdpContent, c.logger)
	default:
		return fmt.Errorf("unsupported codec type: %s", mType)
	}
	if err != nil {
		return fmt.Errorf("failed to create WebM recorder: %w", err)
	}
	defer c.recorder.Close()

	time.Sleep(1 * time.Second)

	// Convert and feed RTP packets
	return c.feedPackets(reader)
}

func (c *RTPDump2WebMConverter) feedPackets(reader *rtpdump.Reader) error {
	startTime := time.Now()

	i := 0
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
		if c.sampleBuilder == nil {
			_ = c.recorder.PushRtpBuf(packet.Payload)
		} else {
			// Unmarshal the RTP packet from the raw payload
			rtpPacket := &rtp.Packet{}
			if err := rtpPacket.Unmarshal(packet.Payload); err != nil {
				c.logger.Warn("Failed to unmarshal RTP packet %d: %v", i, err)
				continue
			}

			// Push packet to samplebuilder for reordering
			c.sampleBuilder.Push(rtpPacket)
		}

		//		time.Sleep(10 * time.Microsecond)
		// Log progress
		if i%10000 == 0 && i > 0 {
			c.logger.Info("Processed %d packets", i)
		}
	}

	if c.sampleBuilder != nil {
		c.sampleBuilder.Flush()
	}

	duration := time.Since(startTime)
	c.logger.Info("Finished feeding %d packets in %v", i, duration)

	// Allow some time for the recorder to finalize
	time.Sleep(2 * time.Second)

	return nil
}

func (c *RTPDump2WebMConverter) buildDefaultReleasePacketHandler() func(pkt *rtp.Packet) {
	return func(pkt *rtp.Packet) {
		if c.lastPkt != nil {
			if pkt.SequenceNumber-c.lastPkt.SequenceNumber > 1 {
				c.logger.Info("Missing Packet Detected, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
			}
		}

		c.lastPkt = pkt

		if e := c.recorder.OnRTP(pkt); e != nil {
			c.logger.Warn("Failed to record RTP packet SeqNum: %d RtpTs: %d: %v", pkt.SequenceNumber, pkt.Timestamp, e)
		}
	}
}

func (c *RTPDump2WebMConverter) buildOpusReleasePacketHandler() func(pkt *rtp.Packet) {
	return func(pkt *rtp.Packet) {
		pkt.SequenceNumber += c.inserted

		if c.lastPkt != nil {
			if pkt.SequenceNumber-c.lastPkt.SequenceNumber > 1 {
				c.logger.Info("Missing Packet Detected, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
			}

			tsDiff := pkt.Timestamp - c.lastPkt.Timestamp // TODO handle rollover
			lastPktDuration := opusPacketDurationMs(c.lastPkt)
			rtpDuration := uint32(lastPktDuration * 48)

			if rtpDuration == 0 {
				rtpDuration = c.lastPktDuration
				c.logger.Info("LastPacket with no duration, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d", c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)
			} else {
				c.lastPktDuration = rtpDuration
			}

			if rtpDuration > 0 && tsDiff > rtpDuration {

				// Calculate how many packets we need to insert, taking care of packet losses
				var toAdd uint16
				if uint32(pkt.SequenceNumber-c.lastPkt.SequenceNumber)*rtpDuration != tsDiff { // TODO handle rollover
					toAdd = uint16(tsDiff/rtpDuration) - (pkt.SequenceNumber - c.lastPkt.SequenceNumber)
				}

				c.logger.Info("Gap detected, inserting %d packets tsDiff %d, Previous SeqNum: %d RtpTs: %d   - Last SeqNum: %d RtpTs: %d",
					toAdd, tsDiff, c.lastPkt.SequenceNumber, c.lastPkt.Timestamp, pkt.SequenceNumber, pkt.Timestamp)

				for i := 1; i <= int(toAdd); i++ {
					ins := c.lastPkt.Clone()
					ins.Payload = ins.Payload[:1] // Keeping only TOC byte
					ins.SequenceNumber += uint16(i)
					ins.Timestamp += uint32(i) * rtpDuration

					c.logger.Debug("Writing inserted Packet %v", ins)
					if e := c.recorder.OnRTP(ins); e != nil {
						c.logger.Warn("Failed to record inserted RTP packet SeqNum: %d RtpTs: %d: %v", ins.SequenceNumber, ins.Timestamp, e)
					}
				}

				c.inserted += toAdd
				pkt.SequenceNumber += toAdd
			}
		}

		c.lastPkt = pkt

		c.logger.Debug("Writing real Packet Last SeqNum: %d RtpTs: %d", pkt.SequenceNumber, pkt.Timestamp)
		if e := c.recorder.OnRTP(pkt); e != nil {
			c.logger.Warn("Failed to record RTP packet SeqNum: %d RtpTs: %d: %v", pkt.SequenceNumber, pkt.Timestamp, e)
		}
	}
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
	var duration int
	switch {
	case config < 3:
		// SILK-only NB: 10, 20, 40 ms
		duration = 100 * (1 << (config & 0x03))
	case config == 3:
		// SILK-only NB: 60 ms
		duration = 600
	case config < 7:
		// SILK-only MB: 10, 20, 40 ms
		duration = 100 * (1 << (config & 0x03))
	case config == 7:
		// SILK-only MB: 60 ms
		duration = 600
	case config < 11:
		// SILK-only WB: 10, 20, 40 ms
		duration = 100 * (1 << (config & 0x03))
	case config == 11:
		// SILK-only WB: 60 ms
		duration = 600
	case config <= 13:
		// Hybrid SWB: 10, 20 ms
		duration = 100 * (1 << (config & 0x01))
	case config <= 15:
		// Hybrid FB: 10, 20 ms
		duration = 100 * (1 << (config & 0x01))
	case config <= 19:
		// CELT-only NB: 2.5, 5, 10, 20 ms
		duration = 25 * (1 << (config & 0x03)) // 2.5ms * 10 for integer math
	case config <= 23:
		// CELT-only WB: 2.5, 5, 10, 20 ms
		duration = 25 * (1 << (config & 0x03)) // 2.5ms * 10 for integer math
	case config <= 27:
		// CELT-only SWB: 2.5, 5, 10, 20 ms
		duration = 25 * (1 << (config & 0x03)) // 2.5ms * 10 for integer math
	case config <= 31:
		// CELT-only FB: 2.5, 5, 10, 20 ms
		duration = 25 * (1 << (config & 0x03)) // 2.5ms * 10 for integer math
	default:
		// MUST NOT HAPPEN
		duration = 0
	}

	frameDuration := float32(duration) / 10

	var frameCount float32
	switch c {
	case 0:
		frameCount = 1
	case 1, 2:
		frameCount = 2
	case 3:
		if len(payload) > 1 {
			frameCount = float32(payload[1] & 0x3F)
		}
	}

	return int(frameDuration * frameCount)
}
