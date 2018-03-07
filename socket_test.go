package websocket

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type CloseableBuffer struct {
	bytes.Buffer
}

func (c *CloseableBuffer) Close() error {
	return nil
}

func TestSend(t *testing.T) {
	ws := &socket{}
	err := ws.Send(0, false, nil)
	if err == nil {
		t.Errorf("Send() should return error when socket state isn't opened.\n")
	}
	ws.state = STATE_OPEN
	var b CloseableBuffer
	ws.conn = &b
	input := make([]byte, 126)
	ws.Send(0, false, input)
	f := frame(b.Bytes())
	if f.getOpcode() != 0 {
		t.Errorf("Send() should be ok\n")
	}
	if f.isFin() {
		t.Errorf("Send() should be ok\n")
	}
}

func TestPing(t *testing.T) {
	ws := &socket{}
	err := ws.Ping(nil)
	if err == nil {
		t.Errorf("Ping() should return error when socket state isn't opened.\n")
	}
	ws.state = STATE_OPEN
	var b CloseableBuffer
	ws.conn = &b
	input := make([]byte, 126)
	ws.Ping(input)
	f := frame(b.Bytes())
	if f.getPayloadLen() != 125 {
		t.Errorf("Ping() should intercept 125 bytes of data.\n")
	}
	if f.getOpcode() != FRAME_TYPE_PING {
		t.Errorf("Ping() should write ping opcode.\n")
	}
	if !f.isFin() {
		t.Errorf("Ping() The control frame should have the fin bit set.\n")
	}
}

func TestPong(t *testing.T) {
	ws := &socket{}
	err := ws.Pong(nil)
	if err == nil {
		t.Errorf("Pong() should return error when socket state isn't opened.\n")
	}
	ws.state = STATE_OPEN
	var b CloseableBuffer
	ws.conn = &b
	input := make([]byte, 126)
	ws.Pong(input)
	f := frame(b.Bytes())
	if f.getPayloadLen() != 125 {
		t.Errorf("Pong() should intercept 125 bytes of data.\n")
	}
	if f.getOpcode() != FRAME_TYPE_PONG {
		t.Errorf("Pong() should write pong opcode.\n")
	}
	if !f.isFin() {
		t.Errorf("Pong() The control frame should have the fin bit set.\n")
	}
}

func TestClose(t *testing.T) {
	ws := &socket{}
	err := ws.Close(0, nil)
	if err == nil {
		t.Errorf("Close() should return error when socket state isn't opened.\n")
	}
	ws.state = STATE_OPEN
	var b CloseableBuffer
	ws.conn = &b
	input := make([]byte, 126)
	ws.Close(0, input)
	f := frame(b.Bytes())
	if binary.BigEndian.Uint16(f[2:]) != uint16(STATUS_NORMAL_CLOSURE) {
		t.Errorf("Close() should set default status code.\n")
	}
	if len(f[2:]) > CONTROL_FRAME_MAX_PAYLOAD_LEN {
		t.Errorf("Close() The control frame max length is %d\n", CONTROL_FRAME_MAX_PAYLOAD_LEN)
	}
	if f.getPayloadLen() != 125 {
		t.Errorf("Close() should intercept 125 bytes of data.\n")
	}
	if f.getOpcode() != FRAME_TYPE_CLOSE {
		t.Errorf("Close() should write close opcode.\n")
	}
	if !f.isFin() {
		t.Errorf("Close() The control frame should have the fin bit set.\n")
	}
}
