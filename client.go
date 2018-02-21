package websocket

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const (
	SEC_WEBSOCKET_KEY_BYTES = 16
	SEC_WEBSOCKET_VERSION   = "13"
)

var (
	conn_err = errors.New("Websocket connection error")
)

// 连接 websocket server
// func Dial(url string) (Conn, error) {

// }

// 客户端握手请求
/*
func handshake(urlstr string) error {
	u, err := url.Parse(urlstr)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	default:
		return errors.New("handshake() url scheme wrong")
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	key := handshakeKey()
	req.Header.Add("Upgrade", "websocket")
	req.Header.Add("Connection", "Upgrade")
	req.Header.Add("Sec-Websocket-Key", key)
	req.Header.Add("Sec-Websocket-Version", SEC_WEBSOCKET_VERSION)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusSwitchingProtocols:
		fmt.Println("ok")
	case http.StatusFound:
		return handshake(resp.Header.Get("Location"))
	default:
		return errors.New("handshake() Connection error: " + resp.Status)
	}

	if !containsInHeader(resp.Header, "Upgrade", "websocket") {
		return conn_err
	}
	if !containsInHeader(resp.Header, "Connection", "Upgrade") {
		return conn_err
	}
	if resp.Header.Get("Sec-Websocket-Accept") != handshakeAccept(key) {
		return conn_err
	}
}
*/

// create 'Sec-Websocket-Key' header value
func handshakeKey() string {
	src := make([]byte, SEC_WEBSOCKET_KEY_BYTES)
	rand.Seed(time.Now().UnixNano())
	for i := range src {
		src[i] = byte(rand.Intn(256))
	}
	return base64.StdEncoding.EncodeToString(src)
}

// Find if the header[key] contains the specified value
func containsInHeader(h http.Header, key, value string) bool {
	for k, v := range h {
		if strings.ToLower(k) == strings.ToLower(key) {
			for _, v := range v {
				if strings.ToLower(v) == strings.ToLower(value) {
					return true
				}
			}
			break
		}
	}
	return false
}
