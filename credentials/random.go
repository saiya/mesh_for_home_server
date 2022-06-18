package credentials

import (
	"crypto/rand"
	"math/big"
)

func NewNonce(letters string) string {
	lettersLen := big.NewInt(int64(len(letters)))
	result := make([]byte, challengeBytes)

	for i := 0; i < challengeBytes; i++ {
		n, err := rand.Int(rand.Reader, lettersLen)
		neverFail(err)
		result[i] = challengeLetters[n.Int64()]
	}
	return string(result)
}
