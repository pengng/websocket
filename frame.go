package websocket

import (
	"encoding/binary"
	"math/rand"
	"time"
)

const (
	FRAME_TYPE_CONTINUATION frameType = iota
	FRAME_TYPE_TEXT
	FRAME_TYPE_BINARY
	FRAME_TYPE_NON_CONTROL1
	FRAME_TYPE_NON_CONTROL2
	FRAME_TYPE_NON_CONTROL3
	FRAME_TYPE_NON_CONTROL4
	FRAME_TYPE_NON_CONTROL5
	FRAME_TYPE_CLOSE
	FRAME_TYPE_PING
	FRAME_TYPE_PONG
	FRAME_TYPE_FURTHER_CONTROL1
	FRAME_TYPE_FURTHER_CONTROL2
	FRAME_TYPE_FURTHER_CONTROL3
	FRAME_TYPE_FURTHER_CONTROL4
	FRAME_TYPE_FURTHER_CONTROL5
)

// 所有控制帧可以带最多125字节数据
const CONTROL_FRAME_MAX_PAYLOAD_LEN = 0x7d

const MAX_PAYLOAD_LEN = ^uint64(0) >> 1

type frameType byte
type frame []byte

func (f frame) fin() frame {
	if len(f) == 0 {
		f = make(frame, 1)
	}
	f[0] |= 0x80
	return f
}

func (f frame) isFin() bool {
	if len(f) == 0 {
		return false
	}
	return f[0]&0x80 != 0
}

func (f frame) mask() frame {
	if len(f) < 2 {
		t := make(frame, 2)
		copy(t, f)
		f = t
	}
	f[1] |= 0x80
	return f
}

func (f frame) isMask() bool {
	if len(f) < 2 {
		return false
	}
	return f[1]&0x80 != 0
}

func (f frame) rsv1() frame {
	if len(f) == 0 {
		f = make(frame, 1)
	}
	f[0] |= 0x40
	return f
}

func (f frame) rsv2() frame {
	if len(f) == 0 {
		f = make(frame, 1)
	}
	f[0] |= 0x20
	return f
}

func (f frame) rsv3() frame {
	if len(f) == 0 {
		f = make(frame, 1)
	}
	f[0] |= 0x10
	return f
}

// opcode 范围为0到15，超出则求余
func (f frame) opcode(t frameType) frame {
	if len(f) == 0 {
		f = make(frame, 1)
	}
	f[0] |= byte(t & 0xf)
	return f
}

func (f frame) getOpcode() frameType {
	if len(f) == 0 {
		return 0
	}
	return frameType(f[0] & 0xff)
}

// The maximum length of the payload is MAX_PAYLOAD_LEN,
// and the excess will be set directly to the maximum value.
func (f frame) payloadLen(bytes uint64) frame {
	if bytes > MAX_PAYLOAD_LEN {
		bytes = MAX_PAYLOAD_LEN
	}
	switch {
	case bytes <= 0x7d:
		t := make(frame, 2)
		copy(t, f)
		f = t
		f[1] |= byte(bytes)
	case bytes < 0x10000:
		t := make(frame, 4)
		copy(t, f)
		f = t
		f[1] |= 0x7e
		binary.BigEndian.PutUint16(f[2:4], uint16(bytes))
	default:
		t := make(frame, 10)
		copy(t, f)
		f = t
		f[1] |= 0x7f
		binary.BigEndian.PutUint64(f[2:10], bytes)
	}
	return f
}

func (f frame) getPayloadLen() uint64 {
	if len(f) < 2 {
		return 0
	}
	l := f[1] & 0x7f
	switch {
	case l == 0x7f:
		return uint64(binary.BigEndian.Uint64(f[2:10]))
	case l == 0x7e:
		return uint64(binary.BigEndian.Uint16(f[2:4]))
	default:
		return uint64(l)
	}
}

func (f frame) maskingKey(key [4]byte) frame {
	l := f.getPayloadLen()
	switch {
	case l > 0xffff:
		t := make(frame, 14)
		copy(t, f)
		f = t
		copy(f[10:14], key[:])
	case l > 0x7d:
		t := make(frame, 8)
		copy(t, f)
		f = t
		copy(f[4:8], key[:])
	default:
		t := make(frame, 6)
		copy(t, f)
		f = t
		copy(f[2:6], key[:])
	}
	return f
}

func createMaskingKey() [4]byte {
	var key [4]byte
	rand.Seed(time.Now().UnixNano())
	for i := range key {
		key[i] = byte(rand.Intn(0x100))
	}
	return key
}

func (f frame) payloadData(d []byte) frame {
	var i int
	l := f.getPayloadLen()
	switch {
	case l > 0xffff:
		i = 10
	case l > 0x7d:
		i = 4
	default:
		i = 2
	}
	if f.isMask() {
		i += 4
	}
	t := make(frame, i+len(d))
	copy(t, f)
	f = t
	copy(f[i:], d)
	return f
}
