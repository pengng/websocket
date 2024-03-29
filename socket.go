package websocket

import (
	"encoding/binary"
	"errors"
	"io"
)

const (
	STATE_CONNECTING readyState = iota
	STATE_OPEN
	STATE_CLOSEING
	STATE_CLOSED
)

const (
	STATUS_NORMAL_CLOSURE statusCode = iota + 1000
	STATUS_GOING_AWAY
	STATUS_PROTOCOL_ERROR
	STATUS_UNACCEPTABLE_MESSAGE_TYPE
	STATUS_RESERVED1
	STATUS_RESERVED2
	STATUS_RESERVED3
	STATUS_INCONSISTENT_MESSAGE_TYPE
	STATUS_MESSAGE_VIOLATES_POLICY
	STATUS_MESSAGE_TOO_BIG
	STATUS_EXPECTED_NEGOTIATE_EXTENSION
	STATUS_UNEXPECTED_CONDITION
	STATUS_FAILURE_TLS_HANDSHAKE statusCode = iota + 1003
)

type statusCode uint16
type readyState byte

type socket struct {
	conn  io.ReadWriteCloser
	state readyState
	mask  bool
	// Get  func() (msgType int, data []byte, err error)
}

func (s *socket) SendText(msg string) error {
	return s.Send(FRAME_TYPE_TEXT, true, []byte(msg))
}

// The frame type range is 0-15, if it is greater than 15 mod 15
func (s *socket) Send(t frameType, fin bool, msg []byte) error {
	if s.state != STATE_OPEN {
		return errors.New("The socket is not opened.")
	}
	var f frame
	if fin {
		f = f.fin()
	}
	f = f.opcode(t).payloadLen(uint64(len(msg)))
	if s.mask {
		k := createMaskingKey()
		f = f.mask().maskingKey(k)
		maskData(k, msg)
	}
	f = f.payloadData(msg)
	_, err := s.conn.Write(f)
	return err
}

// The control frame can be attached with up to 125 bytes of data
// and the overflow part will be discarded.
func (s *socket) Ping(msg []byte) error {
	if len(msg) > CONTROL_FRAME_MAX_PAYLOAD_LEN {
		msg = msg[:CONTROL_FRAME_MAX_PAYLOAD_LEN]
	}
	return s.Send(FRAME_TYPE_PING, true, msg)
}

// The control frame can be attached with up to 125 bytes of data
// and the overflow part will be discarded.
func (s *socket) Pong(msg []byte) error {
	if len(msg) > CONTROL_FRAME_MAX_PAYLOAD_LEN {
		msg = msg[:CONTROL_FRAME_MAX_PAYLOAD_LEN]
	}
	return s.Send(FRAME_TYPE_PONG, true, msg)
}

// Status code can not be less than 1000.
// If less than 1000 is automatically set to 1000.
// The control frame can be attached with up to 125 bytes of data
// and the overflow part will be discarded.
func (s *socket) Close(code statusCode, reason []byte) error {
	if code < STATUS_NORMAL_CLOSURE {
		code = STATUS_NORMAL_CLOSURE
	}
	if len(reason) > CONTROL_FRAME_MAX_PAYLOAD_LEN-2 {
		reason = reason[:CONTROL_FRAME_MAX_PAYLOAD_LEN-2]
	}
	r := make([]byte, len(reason)+2)
	copy(r[2:], reason)
	reason = r
	binary.BigEndian.PutUint16(reason, uint16(code))
	return s.Send(FRAME_TYPE_CLOSE, true, reason)
}
