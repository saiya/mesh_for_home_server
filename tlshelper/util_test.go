package tlshelper

import (
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNeverFailViolated(t *testing.T) {
	assert.Panics(t, func() {
		neverFail(fmt.Errorf("test error"))
	})
}

func TestParseTLSVersion(t *testing.T) {
	assert.Equal(t, uint16(tls.VersionTLS12), parseTLSVersion("tls1.2", 0))
	assert.Equal(t, uint16(tls.VersionTLS12), parseTLSVersion("TLS1.2", 0))

	assert.Equal(t, uint16(tls.VersionTLS13), parseTLSVersion("TLS1.3", 0))

	assert.Equal(t, uint16(tls.VersionTLS13), parseTLSVersion("", tls.VersionTLS13))
}
