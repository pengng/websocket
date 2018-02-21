package websocket

import (
	"crypto/sha1"
	"encoding/base64"
)

const (
	// This substring splices the value of the Sec-Websocket-Key in the client's request header to generate Sec-Websocket-Accept that responds to headers
	SUBSTRING_FOR_GENERATE_HANDSHAKE_ACCEPT = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

func Listen() {

}

func Upgrade() {

}

func handshakeAccept(key string) string {
	sum := sha1.Sum([]byte(key + SUBSTRING_FOR_GENERATE_HANDSHAKE_ACCEPT))
	return base64.StdEncoding.EncodeToString(sum[:])
}
