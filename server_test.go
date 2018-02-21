package websocket

import "testing"

func TestHandshakeAccept(t *testing.T) {
	input := []string{
		"dGhlIHNhbXBsZSBub25jZQ==",
		"gH0TL3Qv3jdMqmasceVkUg==",
		"TLqKbCEUW8GMCnlA3b0nEg=="}
	expected := []string{
		"s3pPLMBiTxaQ9kYGzzhZRbK+xOo=",
		"7KSHspfnUahkAsYWUU9imQu9eB4=",
		"nyuhCaUJfpGB+uOWXL0JDwokzlg="}
	for i, v := range input {
		if handshakeAccept(v) != expected[i] {
			t.Errorf("handshakeAccept(%q) should return %q", v, expected[i])
		}
	}
}

func BenchmarkHandshakeAccept(b *testing.B) {
	for i := 0; i < b.N; i++ {
		handshakeAccept("dGhlIHNhbXBsZSBub25jZQ==")
	}
}
