package testutil

import "math/rand"

func RandomBytes(len int) []byte {
	buf := make([]byte, len)
	rand.Read(buf)
	return buf
}
