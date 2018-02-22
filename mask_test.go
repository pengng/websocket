package websocket

import "testing"

func TestMaskData(t *testing.T) {
	d := []byte{1, 2, 3, 4, 5, 6, 7}
	k := [4]byte{4, 3, 2, 1}
	expected := []byte{5, 1, 1, 5, 1, 5, 5}
	maskData(k, d)
	for i, v := range expected {
		if d[i] != v {
			t.Errorf("maskData() should be ok")
		}
	}
}

func BenchmarkMaskData(b *testing.B) {
	d := []byte{1, 2, 3, 4, 5, 6, 7}
	k := [4]byte{4, 3, 2, 1}
	for i := 0; i < b.N; i++ {
		maskData(k, d)
	}
}
