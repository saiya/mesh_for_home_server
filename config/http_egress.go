package config

import (
	"net"
	"net/http"
	"time"
)

type HTTPEgressConfig struct {
	// (optional) This egress receive HTTP traffic that match with this hostname.
	// By default or for empty string, match with any host.
	//
	// Can use "*" in the 1st part of hostname (e.g. "*.example.com" match with "sub.example.com").
	// "*" match with only 1 level of hostname element (e.g. "*.example.com" not match with "a.b.example.com")
	//
	// If multiple egress matches with request, longest match preffered.
	Host string `json:"host" yaml:"host"`

	// Upstream serer, "schme://hostname:port" format (e.g. "http://localhost:8080")
	Server string `json:"server" yaml:"server"`

	MaxConnections   int          `json:"max_connections" yaml:"max_connections"`
	KeepAliveTimeout string       `json:"keep_alive_timeout" yaml:"keep_alive_timeout"`
	ConnectTimeout   string       `json:"connect_timeout" yaml:"connect_timeout"`
	ResponseTimeout  *HTTPTimeout `json:"response_timeout" yaml:"response_timeout"`
}

const HTTPEgressMaxConnectionsDefault = 100
const HTTPEgressKeepAliveTimeoutDefault = 60 * time.Second
const HTTPEgressConnectTimeoutDefault = 10 * time.Second

func (c *HTTPEgressConfig) ConfigureDialer(d *net.Dialer) {
	d.Timeout = parseDurationOrDefault(c.ConnectTimeout, HTTPEgressConnectTimeoutDefault)
}

func (c *HTTPEgressConfig) ConfigureHTTPTransport(t *http.Transport) {
	maxConns := HTTPEgressMaxConnectionsDefault
	if c.MaxConnections != 0 {
		maxConns = c.MaxConnections
	}
	t.MaxIdleConns = maxConns
	t.MaxIdleConnsPerHost = maxConns
	t.MaxConnsPerHost = maxConns

	t.IdleConnTimeout = parseDurationOrDefault(c.KeepAliveTimeout, HTTPEgressKeepAliveTimeoutDefault)
	t.TLSHandshakeTimeout = parseDurationOrDefault(c.ConnectTimeout, HTTPEgressConnectTimeoutDefault)
	t.ResponseHeaderTimeout = c.ResponseTimeout.HeaderTimeout()
	t.ExpectContinueTimeout = c.ResponseTimeout.HeaderTimeout()
}
