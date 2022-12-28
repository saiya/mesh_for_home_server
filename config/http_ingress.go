package config

import "time"

type HTTPIngressConfig struct {
	// (optional) "host:port" to listen (e.g. "localhost:8080", ":8080")
	// Use random available port by default.
	Listen string `json:"listen" yaml:"listen"`

	TLS *TLSServerConfig `json:"tls" yaml:"tls"`

	Probe *HTTPProbeConfig `json:"probe" yaml:"probe"`

	RequestTimeout   *HTTPTimeout `json:"request_timeout" yaml:"request_timeout"`
	ResponseTimeout  *HTTPTimeout `json:"response_timeout" yaml:"response_timeout"`
	KeepAliveTimeout string       `json:"keep_alive_timeout" yaml:"keep_alive_timeout"`
}

const HTTPIngressKeepAliveTimeoutDefault = 60 * time.Second

func (c *HTTPIngressConfig) GetKeepAliveTimeout() time.Duration {
	return parseDurationOrDefault(c.KeepAliveTimeout, HTTPIngressKeepAliveTimeoutDefault)
}
