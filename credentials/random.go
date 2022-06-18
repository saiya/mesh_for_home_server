package credentials

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func NewNonce(letters string) string {
	lettersLen := big.NewInt(int64(len(letters)))
	result := make([]byte, challengeBytes)

	for i := 0; i < challengeBytes; i++ {
		n, err := rand.Int(rand.Reader, lettersLen)
		if err != nil {
			panic(fmt.Errorf("failed to generate rand: %w", err))
		}
		result[i] = challengeLetters[n.Int64()]
	}
	return string(result)
}
