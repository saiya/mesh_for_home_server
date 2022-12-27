package config

import "time"

type HTTPTimeout struct {
	Header string `json:"header" yaml:"header"`
	Body   string `json:"body" yaml:"body"`
}

var HTTPHeaderTimeoutDefault = 60 * time.Second
var HTTPBodyTimeoutDefault = 60 * time.Second

func (ht *HTTPTimeout) HeaderTimeout() time.Duration {
	if ht == nil {
		return HTTPHeaderTimeoutDefault
	}
	return parseDurationOrDefault(ht.Header, HTTPHeaderTimeoutDefault)
}

func (ht *HTTPTimeout) BodyTimeout() time.Duration {
	if ht == nil {
		return HTTPBodyTimeoutDefault
	}
	return parseDurationOrDefault(ht.Body, HTTPBodyTimeoutDefault)
}

type HTTPProbeConfig struct {
	// (optional) Hostname (expected Host header value) of the probe endpoint
	Host string `json:"host" yaml:"host"`
	// Path of the probe endpoint (e.g. "/probe")
	Path string `json:"path" yaml:"path"`
}
