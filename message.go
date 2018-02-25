package websocket

import (
	"encoding/binary"
)

type closeMsg []byte

// Status code can not be less than 1000.
// If less than 1000 is automatically set to 1000.
// The control frame can be attached with up to 125 bytes of data
// and the overflow part will be discarded.
func NewCloseMsg(code statusCode, msg string) closeMsg {
	if code < STATUS_NORMAL_CLOSURE {
		code = STATUS_NORMAL_CLOSURE
	}
	if len(msg) > CONTROL_FRAME_MAX_PAYLOAD_LEN-2 {
		msg = msg[:CONTROL_FRAME_MAX_PAYLOAD_LEN-2]
	}
	m := make(closeMsg, len(msg)+2)
	copy(m[2:], []byte(msg))
	binary.BigEndian.PutUint16(m, uint16(code))
	return m
}
