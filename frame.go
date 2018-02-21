package websocket

type frame []byte

func (f frame) fin() {
	f[0] |= 1 << 7
}

func (f frame) fin1() {
	f[0] |= 0x80
}

func (f frame) isFin() bool {
	return f[0]&1<<7 != 0
}

func (f frame) mask() {

}

func (f frame) rsv1() {

}

func (f frame) rsv2() {

}

func (f frame) rsv3() {

}

func (f frame) opcode(c byte) {

}

func (f frame) payloadLen(l int) {

}

func (f frame) maskingKey(key [32]byte) {

}

func (f frame) payloadData(data []byte) {

}
