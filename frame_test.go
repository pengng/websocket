package websocket

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestFin(t *testing.T) {
	var f frame
	f = f.fin()
	if f[0]&0x80 == 0 {
		t.Errorf("fin() should set first bit to 1")
	}
}

func TestIsFin(t *testing.T) {
	var f frame
	if f.isFin() {
		t.Errorf("isFin() should return false")
	}
	f = f.fin()
	if !f.isFin() {
		t.Errorf("isFin() should return true")
	}
}

func TestRsv1(t *testing.T) {
	var f frame
	f = f.rsv1()
	if f[0]&0x40 == 0 {
		t.Errorf("rsv1() should set second bit to 1")
	}
}

func TestRsv2(t *testing.T) {
	var f frame
	f = f.rsv2()
	if f[0]&0x20 == 0 {
		t.Errorf("rsv2() should set third bit to 1")
	}
}

func TestRsv3(t *testing.T) {
	var f frame
	f = f.rsv3()
	if f[0]&0x10 == 0 {
		t.Errorf("rsv3() should set four bit to 1")
	}
}

func TestOpcode(t *testing.T) {
	var f frame
	f = f.opcode(FRAME_TYPE_TEXT)
	if f[0]&(byte(FRAME_TYPE_TEXT)) == 0 {
		t.Errorf("opcode should be ok")
	}
}

func TestMask(t *testing.T) {
	var f frame
	f = f.mask()
	if f[1]&0x80 == 0 {
		t.Errorf("mask() should set 8th bit to 1")
	}
	f = make(frame, 1)
	f = f.mask()
	if f[1]&0x80 == 0 {
		t.Errorf("mask() should set 8th bit to 1")
	}
	f = make(frame, 2)
	f = f.mask()
	if f[1]&0x80 == 0 {
		t.Errorf("mask() should set 8th bit to 1")
	}
}

func TestIsMask(t *testing.T) {
	var f frame
	if f.isMask() {
		t.Errorf("isMask() should return false")
	}
	f = f.mask()
	if !f.isMask() {
		t.Errorf("isMask() should return true")
	}
}

func TestPayloadLen(t *testing.T) {
	var f frame
	var input uint64 = 0x7d
	f = f.payloadLen(input)
	if f[1]&0x7f != byte(input) {
		t.Errorf("payloadLen(%#x) should set 2th byte to %#[1]x", input)
	}

	f = frame{}
	input = 0x7e
	f = f.payloadLen(input)
	if !(f[1]&0x7f == 0x7e && binary.BigEndian.Uint16(f[2:4]) == uint16(input)) {
		t.Errorf("payloadLen(%#x) should set 2th byte to 0x7e\nset 3-4th byte to %#[1]x", input)
	}

	f = frame{}
	input = 0x10000
	f = f.payloadLen(input)
	if !(f[1]&0x7f == 0x7f && binary.BigEndian.Uint64(f[2:10]) == input) {
		t.Errorf("payloadLen(%#x) should set 2th byte to 0x7f\nset 3-10th byte to %#[1]x", input)
	}
}

func TestGetPayloadLen(t *testing.T) {
	input := []uint64{0x00, 0x7d, 0x7e, 0x10000}
	for _, v := range input {
		var f frame
		f = f.payloadLen(v)
		if f.getPayloadLen() != v {
			t.Errorf("getPayloadLen() should return %#x", v)
		}
	}
}

func TestMaskingKey(t *testing.T) {
	var f frame
	k := createMaskingKey()
	f = f.maskingKey(k)
	for i, v := range k {
		if f[2+i] != v {
			t.Errorf("maskingKey() should be ok")
		}
	}

	f = frame{}
	f = f.payloadLen(0x7e)
	f = f.maskingKey(k)
	for i, v := range k {
		if f[4+i] != v {
			t.Errorf("maskingKey() should be ok")
		}
	}

	f = frame{}
	f = f.payloadLen(0x10000)
	f = f.maskingKey(k)
	for i, v := range k {
		if f[10+i] != v {
			t.Errorf("maskingKey() should be ok")
		}
	}
}

func TestCreateMaskingKey(t *testing.T) {
	x, y := createMaskingKey(), createMaskingKey()
	if x == y {
		t.Errorf("createMaskingKey() should return random [4]byte")
	}
	if reflect.TypeOf(x).String() != "[4]uint8" {
		t.Errorf("createMaskingKey() should return type [4]byte")
	}
}

func TestPayloadData(t *testing.T) {
	var f frame
	d := []byte("123456789")
	f = f.payloadData(d)
	for i, v := range d {
		if f[2+i] != v {
			t.Errorf("payloadData() should be ok")
		}
	}

	f = f.mask()
	f = f.payloadData(d)
	for i, v := range d {
		if f[6+i] != v {
			t.Errorf("payloadData() should be ok")
		}
	}

	f = frame{}
	f = f.mask()
	f = f.payloadLen(0x10000)
	f = f.payloadData(d)
	for i, v := range d {
		if f[14+i] != v {
			t.Errorf("payloadData() should be ok")
		}
	}
}

func BenchmarkCreateMaskingKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		createMaskingKey()
	}
}
