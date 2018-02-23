package websocket

import (
	"net"
)

const (
	OPEN connState = iota
	CLOSEING
	CLOSED
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
type connState byte

type socket struct {
	conn    net.Conn
	GetMsg  func() (msgType int, data []byte, err error)
	SendMsg func(msgType int, data interface{}) (err error)
	Close   func() error
}
