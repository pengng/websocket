package websocket

import (
	"bytes"
	"testing"
)

type CloseableBuffer struct {
	bytes.Buffer
}

func (c *CloseableBuffer) Close() error {
	return nil
}

func TestPing(t *testing.T) {
	ws := &socket{}
	err := ws.Ping(nil)
	if err == nil {
		t.Errorf("Ping() should return error when socket state isn't opened.")
	}
	ws.state = STATE_OPEN
	var b CloseableBuffer
	ws.conn = &b
	input := make([]byte, 126)
	ws.Ping(input)
	f := frame(b.Bytes())
	if f.getPayloadLen() != 125 {
		t.Errorf("Ping() should intercept 125 bytes of data")
	}
	if f.getOpcode() != FRAME_TYPE_PING {
		t.Errorf("Ping() should write ping opcode.")
	}
}
