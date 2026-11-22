// SPDX-FileCopyrightText: 2023 The Pion community <https://pion.ly>
// SPDX-License-Identifier: MIT

package rtpreorder

type sampleSequenceLocation struct {
	head uint16
	tail uint16
}

func (l sampleSequenceLocation) empty() bool {
	return l.head == l.tail
}

func (l sampleSequenceLocation) hasData() bool {
	return l.head != l.tail
}

func (l sampleSequenceLocation) count() uint16 {
	return seqnumDistance(l.head, l.tail)
}

const (
	slCompareVoid = iota
	slCompareBefore
	slCompareInside
	slCompareAfter
)

func (l sampleSequenceLocation) compare(pos uint16) int {
	if l.head == l.tail {
		return slCompareVoid
	}

	if l.head < l.tail {
		if l.head <= pos && pos < l.tail {
			return slCompareInside
		}
	} else {
		if l.head <= pos || pos < l.tail {
			return slCompareInside
		}
	}

	if l.head-pos <= pos-l.tail {
		return slCompareBefore
	}

	return slCompareAfter
}
