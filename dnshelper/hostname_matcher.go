package dnshelper

import (
	"regexp"
	"strings"
)

// HostnameMatcher generate matcher function
// Matcher returns positive int or 0 represents priority if matched.
// If not match, matcher returns minus value.
func HostnameMatcher(matcher string) func(string) int64 {
	host := normalizeHostname(matcher)

	if host == "" {
		return func(givenHostname string) int64 {
			return 0
		}
	} else if strings.HasPrefix(host, "*.") || host == "*" {
		hostSuffix := host[1:]
		regex := regexp.MustCompile("^[^.]+" + regexp.QuoteMeta(hostSuffix) + "$")
		return func(givenHostname string) int64 {
			if regex.MatchString(normalizeHostname(givenHostname)) {
				return int64(len(host))
			}
			return -1
		}
	} else {
		return func(givenHostname string) int64 {
			if normalizeHostname(givenHostname) == host {
				return int64(len(host))
			}
			return -1
		}
	}
}

func normalizeHostname(s string) string {
	return strings.TrimSuffix(strings.ToLower(s), ".")
}
