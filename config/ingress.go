package config

type IngressConfig struct {
	HTTP []HTTPIngressConfig `json:"http"`
}

type HTTPIngressConfig struct {
	// (optional) "host:port" to listen (e.g. "localhost:8080", ":8080")
	// Use random available port by default.
	Listen string `json:"listen"`

	TLS *TLSServerConfig `json:"tls"`

	Probe *HTTPProbeConfig `json:"probe"`
}

type HTTPProbeConfig struct {
	// (optional) Hostname (expected Host header value) of the probe endpoint
	Host string
	// Path of the probe endpoint (e.g. "/probe")
	Path string
}
