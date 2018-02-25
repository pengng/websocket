package websocket

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
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
func Dial(rawurl string, config *tls.Config) (*socket, error) {

	u, err := parseUrl(rawurl)
	if err != nil {
		return nil, err
	}
	var conn net.Conn
	if u.Scheme == "ws" {
		conn, err = net.Dial("tcp", u.Host)
	} else {
		conn, err = tls.Dial("tcp", u.Host, config)
	}
	if err != nil {
		return nil, err
	}
	k := handshakeKey()
	err = handshake(conn, u, k)
	if err != nil {
		return nil, err
	}
	err = parseHandshake(conn, k)
	if err != nil {
		return nil, err
	}
	return &socket{conn: conn, state: STATE_OPEN, mask: true}, nil
}

func parseUrl(rawurl string) (*url.URL, error) {
	if ok, _ := regexp.MatchString(`^\w+://\w+(?::\d{1, 5})?`, rawurl); !ok {
		return nil, errors.New(fmt.Sprintf("parseUrl() The URL %q is invalid.", rawurl))
	}
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}
	if !(u.Scheme == "ws" || u.Scheme == "wss") {
		return nil, errors.New(fmt.Sprintf("parseUrl() The URL's scheme must be either 'ws' or 'wss'. %q is not allowed.", u.Scheme))
	}
	if u.Port() == "" {
		if u.Scheme == "ws" {
			u.Host += ":80"
		} else {
			u.Host += ":443"
		}
	}
	if u.Path == "" {
		u.Path = "/"
	}
	return u, nil
}

func handshake(w io.Writer, u *url.URL, key string) error {
	if ok, _ := regexp.MatchString(`^[a-zA-Z0-9+/]{22}==$`, key); !ok {
		return errors.New("handshake() The key must be base64-encoded 16-bit bytes")
	}
	fmt.Fprintf(w, "GET %s HTTP/1.1\r\n", u.Path)
	h := make(http.Header)
	h.Add("Host", u.Host)
	h.Add("Upgrade", "websocket")
	h.Add("Connection", "Upgrade")
	h.Add("Sec-Websocket-Key", key)
	h.Add("Sec-Websocket-Version", SEC_WEBSOCKET_VERSION)
	h.Add("Content-Length", "0")
	h.Write(w)
	_, err := fmt.Fprint(w, "\r\n")
	return err
}

func parseHandshake(r io.Reader, key string) error {
	rd := bufio.NewScanner(r)
	if !rd.Scan() {
		if err := rd.Err(); err == io.EOF {
			return errors.New("parseHandshake() The server not headline response.")
		} else {
			return err
		}
	}
	l := rd.Text()
	a := strings.Split(l, " ")
	if len(a) < 3 {
		return errors.New(fmt.Sprintf("parseHandshake() The text %q don't like headline.", l))
	}
	c, err := strconv.ParseUint(a[1], 10, 8)
	if err != nil {
		return err
	}
	h := make(http.Header)
	for rd.Scan() {
		s := strings.Split(rd.Text(), ": ")
		if len(s) == 2 {
			h.Add(s[0], s[1])
		} else {
			break
		}
	}
	if c != http.StatusSwitchingProtocols {
		return errors.New("parseHandshake() Connection error")
	}

	if !containsInHeader(h, "Upgrade", "websocket") {
		return conn_err
	}
	if !containsInHeader(h, "Connection", "Upgrade") {
		return conn_err
	}
	if h.Get("Sec-Websocket-Accept") != handshakeAccept(key) {
		return conn_err
	}
	return nil
}

// create 'Sec-Websocket-Key' header value
func handshakeKey() string {
	src := make([]byte, SEC_WEBSOCKET_KEY_BYTES)
	rand.Seed(time.Now().UnixNano())
	for i := range src {
		src[i] = byte(rand.Intn(0x100))
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
