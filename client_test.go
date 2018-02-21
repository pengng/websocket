package websocket

import (
	"net/http"
	"testing"
)

const (
	SEC_WEBSOCKET_KEY_BASE64_BYTES = 24
)

func TestHandshake(t *testing.T) {
	a, b := handshakeKey(), handshakeKey()
	if !(len(a) == SEC_WEBSOCKET_KEY_BASE64_BYTES && len(a) == len(b)) {
		t.Errorf("handshake() Sec-Websocket-Key should be 16 bytes")
	}
	if a == b {
		t.Errorf("handshake() Sec-Websocket-Key should be random")
	}
}

func BenchmarkHandshake(b *testing.B) {
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
