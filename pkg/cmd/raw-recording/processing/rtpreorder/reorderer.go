// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

// Package rtpreorder buffers RTP packets and releases them in sample order.
//
// It is based on github.com/pion/webrtc/v4/pkg/media/samplebuilder, but it
// deliberately avoids building media.Sample values. The raw recorder only needs
// codec-aware packet reordering before forwarding RTP packets to GStreamer; it
// does not need decoded sample payloads, durations, metadata, or RTP headers.
package rtpreorder

import (
	"math"
	"time"

	"github.com/pion/rtp"
)

// Reorderer buffers packets until media frames are complete or late enough to
// release. maxLate is measured in RTP sequence numbers.
type Reorderer struct {
	maxLate          uint16
	maxLateTimestamp uint32
	buffer           [math.MaxUint16 + 1]*rtp.Packet

	depacketizer rtp.Depacketizer
	sampleRate   uint32

	packetReleaseHandler func(*rtp.Packet)

	filled sampleSequenceLocation
	active sampleSequenceLocation

	lastSampleTimestamp *uint32

	droppedPackets uint16
	paddingPackets uint16
}

// New constructs a Reorderer.
func New(maxLate uint16, depacketizer rtp.Depacketizer, sampleRate uint32, opts ...Option) *Reorderer {
	r := &Reorderer{
		maxLate:      maxLate,
		depacketizer: depacketizer,
		sampleRate:   sampleRate,
	}
	for _, o := range opts {
		o(r)
	}

	return r
}

func (r *Reorderer) tooOld(location sampleSequenceLocation) bool {
	if r.maxLateTimestamp == 0 {
		return false
	}

	var foundHead *rtp.Packet
	var foundTail *rtp.Packet

	for i := location.head; i != location.tail; i++ {
		if packet := r.buffer[i]; packet != nil {
			foundHead = packet

			break
		}
	}

	if foundHead == nil {
		return false
	}

	for i := location.tail - 1; i != location.head; i-- {
		if packet := r.buffer[i]; packet != nil {
			foundTail = packet

			break
		}
	}

	if foundTail == nil {
		return false
	}

	return timestampDistance(foundHead.Timestamp, foundTail.Timestamp) > r.maxLateTimestamp
}

func (r *Reorderer) fetchTimestamp(location sampleSequenceLocation) (timestamp uint32, hasData bool) {
	if location.empty() {
		return 0, false
	}
	packet := r.buffer[location.head]
	if packet == nil {
		return 0, false
	}

	return packet.Timestamp, true
}

func (r *Reorderer) releasePacket(i uint16) {
	var packet *rtp.Packet
	packet, r.buffer[i] = r.buffer[i], nil
	if packet != nil && r.packetReleaseHandler != nil {
		r.packetReleaseHandler(packet)
	}
}

func (r *Reorderer) purgeConsumedBuffers() {
	r.purgeConsumedLocation(r.active, false)
}

func (r *Reorderer) purgeConsumedLocation(consume sampleSequenceLocation, forceConsume bool) {
	if !r.filled.hasData() {
		return
	}

	switch consume.compare(r.filled.head) {
	case slCompareInside:
		if !forceConsume {
			break
		}

		fallthrough
	case slCompareBefore:
		r.releasePacket(r.filled.head)
		r.filled.head++
	}
}

func (r *Reorderer) purgeBuffers(flush bool) {
	r.purgeConsumedBuffers()

	for (r.tooOld(r.filled) || r.filled.count() > r.maxLate || flush) && r.filled.hasData() {
		if r.active.empty() {
			r.active = r.filled
		}

		if r.active.hasData() && r.active.head == r.filled.head {
			if r.releaseSample(true) {
				continue
			}

			r.active.head++
			r.droppedPackets++
		}

		r.releasePacket(r.filled.head)
		r.filled.head++
	}
}

// Push adds an RTP packet to the reorder buffer. Push does not copy the input.
func (r *Reorderer) Push(packet *rtp.Packet) {
	r.buffer[packet.SequenceNumber] = packet

	switch r.filled.compare(packet.SequenceNumber) {
	case slCompareVoid:
		r.filled.head = packet.SequenceNumber
		r.filled.tail = packet.SequenceNumber + 1
	case slCompareBefore:
		r.filled.head = packet.SequenceNumber
	case slCompareAfter:
		r.filled.tail = packet.SequenceNumber + 1
	case slCompareInside:
		break
	}
	r.purgeBuffers(false)
}

// Flush releases all packets that can be released from the buffer.
func (r *Reorderer) Flush() {
	r.purgeBuffers(true)
}

func (r *Reorderer) releaseSample(purgingBuffers bool) bool {
	if r.active.empty() {
		r.active = r.filled
	}

	if r.active.empty() {
		return false
	}

	if r.filled.compare(r.active.tail) == slCompareInside {
		r.active.tail = r.filled.tail
	}

	var consume sampleSequenceLocation

	for i := r.active.head; r.buffer[i] != nil && r.active.compare(i) != slCompareAfter; i++ {
		if r.depacketizer.IsPartitionTail(r.buffer[i].Marker, r.buffer[i].Payload) {
			consume.head = r.active.head
			consume.tail = i + 1

			break
		}
		headTimestamp, hasData := r.fetchTimestamp(r.active)
		if hasData && r.buffer[i].Timestamp != headTimestamp {
			consume.head = r.active.head
			consume.tail = i

			break
		}
	}

	if consume.empty() {
		return false
	}

	if !purgingBuffers && r.buffer[consume.tail] == nil {
		return false
	}

	sampleTimestamp, _ := r.fetchTimestamp(r.active)

	r.active.head = consume.tail

	if !r.depacketizer.IsPartitionHead(r.buffer[consume.head].Payload) {
		isPadding := false
		for i := consume.head; i != consume.tail; i++ {
			if r.lastSampleTimestamp != nil && *r.lastSampleTimestamp == r.buffer[i].Timestamp && len(r.buffer[i].Payload) == 0 {
				isPadding = true
			}
		}
		r.droppedPackets += consume.count()
		if isPadding {
			r.paddingPackets += consume.count()
		}
		r.purgeConsumedLocation(consume, true)
		r.purgeConsumedBuffers()

		return false
	}

	r.droppedPackets = 0
	r.paddingPackets = 0
	r.lastSampleTimestamp = new(uint32)
	*r.lastSampleTimestamp = sampleTimestamp

	r.purgeConsumedLocation(consume, true)
	r.purgeConsumedBuffers()

	return true
}

func seqnumDistance(x, y uint16) uint16 {
	diff := int16(x - y) //nolint:gosec // G115
	if diff < 0 {
		return uint16(-diff)
	}

	return uint16(diff)
}

func timestampDistance(x, y uint32) uint32 {
	diff := int32(x - y) //nolint:gosec // G115
	if diff < 0 {
		return uint32(-diff)
	}

	return uint32(diff)
}

// Option configures a Reorderer.
type Option func(*Reorderer)

// WithPacketReleaseHandler sets a callback when the reorderer removes its
// reference to a packet.
func WithPacketReleaseHandler(h func(*rtp.Packet)) Option {
	return func(r *Reorderer) {
		r.packetReleaseHandler = h
	}
}

// WithMaxTimeDelay ensures packets that are too old in the buffer get purged
// based on timestamp distance rather than sequence distance alone.
func WithMaxTimeDelay(maxLateDuration time.Duration) Option {
	return func(r *Reorderer) {
		totalMillis := maxLateDuration.Milliseconds()
		r.maxLateTimestamp = uint32(int64(r.sampleRate) * totalMillis / 1000) //nolint:gosec // G115
	}
}
