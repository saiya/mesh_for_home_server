package auth

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNeverFailViolated(t *testing.T) {
	assert.Panics(t, func() {
		neverFail(fmt.Errorf("test error"))
	})
}
