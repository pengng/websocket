package websocket

func maskData(k [4]byte, d []byte) {
	for i := range d {
		d[i] ^= k[i%4]
	}
}
