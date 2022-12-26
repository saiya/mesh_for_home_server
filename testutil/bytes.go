package testutil

import "math/rand"

func RandomBytes(len int) []byte {
	buf := make([]byte, len)
	_, err := rand.Read(buf)
	if err != nil {
		panic("failed to generate random")
	}
	return buf
}
