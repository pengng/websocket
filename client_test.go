package websocket

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

const (
	SEC_WEBSOCKET_KEY_BASE64_BYTES = 24
)

func TestHandshakeKey(t *testing.T) {
	a, b := handshakeKey(), handshakeKey()
	if !(len(a) == SEC_WEBSOCKET_KEY_BASE64_BYTES && len(a) == len(b)) {
		t.Errorf("handshakeKey() Sec-Websocket-Key should be %d bytes", SEC_WEBSOCKET_KEY_BYTES)
	}
	if a == b {
		t.Errorf("handshakeKey() Sec-Websocket-Key should be random")
	}
}

func TestParseUrl(t *testing.T) {
	in := []string{"", "1234", "ss://localhost:3000"}
	for _, v := range in {
		if _, err := parseUrl(v); err == nil {
			t.Errorf("parseUrl(%q) should return err", in)
		}
	}
}

func TestHandshake(t *testing.T) {
	var b bytes.Buffer
	var p, k = "/chat", handshakeKey()
	err := handshake(&b, p, k)
	if err != nil {
		t.Errorf("handshake() should be ok")
	}
	s := fmt.Sprintf("GET %s HTTP/1.1\r\n\r\nConnection: Upgrade\r\nSec-Websocket-Key: %s\r\nSec-Websocket-Version: 13\r\nUpgrade: websocket\r\n\r\n", p, k)
	if b.String() != s {
		t.Errorf("handshake() The correct http message should be written")
	}
	k = "111111111111111111111111"
	err = handshake(&b, p, k)
	if err == nil {
		t.Errorf("handshake() should return err")
	}
}

func BenchmarkHandshakeKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handshakeKey()
	}
}

func TestContainsInHeader(t *testing.T) {
	h := http.Header{
		"otherHeader": []string{"value"},
		"Upgrade":     []string{"websocket"}}

	if !containsInHeader(h, "Upgrade", "websocket") {
		t.Errorf("containsInHeader() should be ok")
	}
	if containsInHeader(h, "test", "value") {
		t.Errorf("containsInHeader() should return false if not find key")
	}
	if containsInHeader(h, "otherHeader", "val") {
		t.Errorf("containsInHeader() should return false if find key but not find value")
	}
	if !containsInHeader(h, "UpGraDe", "weBSocket") {
		t.Errorf("containsInHeader() should be ignore case")
	}
}

func BenchmarkContainsInHeader(b *testing.B) {
	h := http.Header{
		"otherHeader": []string{"value"},
		"Upgrade":     []string{"websocket"}}

	for i := 0; i < b.N; i++ {
		containsInHeader(h, "Upgrade", "websocket")
	}
}
