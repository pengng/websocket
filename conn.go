package websocket

type Conn interface {
	GetMsg() (msgType int, data []byte, err error)

	SendMsg(msgType int, data interface{}) (err error)

	Close() error
}
