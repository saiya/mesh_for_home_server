package tlshelper

import (
	"crypto/tls"
	"strings"
)

// neverFail panics program only if non-nil error given
func neverFail(err error) {
	if err == nil {
		return
	}

	panic(err)
}

func parseTLSVersion(s string, defaultValue uint16) uint16 {
	switch strings.ToUpper(s) {
	case "TLS1.2":
		return tls.VersionTLS12
	case "TLS1.3":
		return tls.VersionTLS13
	default:
		return defaultValue
	}
}
