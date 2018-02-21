package websocket

import "testing"

var f frame

func BenchmarkFin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f.fin()
	}
}

func BenchmarkFin1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		f.fin1()
	}
}
