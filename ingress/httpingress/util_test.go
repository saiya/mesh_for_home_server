package httpingress

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

func TestWarnErr(t *testing.T) {
	warnIfError("test error", fmt.Errorf("test error"))
}

func TestDebugErr(t *testing.T) {
	debugIfError("test error", fmt.Errorf("test error"))
}
