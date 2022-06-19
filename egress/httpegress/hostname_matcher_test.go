package httpegress

import (
	"testing"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/stretchr/testify/assert"
)

func TestHostnameMatcher(t *testing.T) {
	test := func(hostnameMatchWith string, requestHostname string, expect bool) {
		matcher := hostnameMatcher(
			&config.HTTPEgressConfig{
				Host: hostnameMatchWith,
			},
		)
		assert.Equal(t, expect, matcher(requestHostname) >= 0)
	}

	test("", "example.com", true)
	test("*", "example.com", false) // "*" match with only 1-level of hostname element
	test("*.com", "example.com", true)
	test("example.com", "example.com", true)
	test("example.com", "sub.example.com", false)
	test("*.example.com", "sub.example.com", true)
	test("*.example.com", "sub.sub.example.com", false)
	test("*.sub.example.com", "sub.example.com", false)
}

func TestHostnameNormalization(t *testing.T) {
	assert.Equal(t, "", normalizeHostname(""))
	assert.Equal(t, "", normalizeHostname("."))
	assert.Equal(t, "*", normalizeHostname("*"))

	assert.Equal(t, "example.com", normalizeHostname("example.com"))
	assert.Equal(t, "example.com", normalizeHostname("example.com."))
	assert.Equal(t, "*.example.com", normalizeHostname("*.example.com"))
}
