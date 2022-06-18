package credentials

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecureRandomString(t *testing.T) {
	const letters = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	assert.NotEqual(t, SecureRandomString(32, letters), SecureRandomString(32, letters))

	assert.Equal(t, 17, len(SecureRandomString(17, letters)))
	assert.Regexp(t, "^["+letters+"]+$", SecureRandomString(32, letters))
}
