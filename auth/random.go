package auth

import (
	"crypto/rand"
	"math/big"
)

func SecureRandomString(length int, letters string) string {
	lettersLen := big.NewInt(int64(len(letters)))
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, lettersLen)
		neverFail(err)
		result[i] = challengeLetters[n.Int64()]
	}
	return string(result)
}
