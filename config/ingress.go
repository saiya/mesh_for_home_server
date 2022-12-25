package config

type IngressConfigs struct {
	HTTP []HTTPIngressConfig `json:"http" yaml:"http"`
}

type HTTPIngressConfig struct {
	// (optional) "host:port" to listen (e.g. "localhost:8080", ":8080")
	// Use random available port by default.
	Listen string `json:"listen" yaml:"listen"`

	TLS *TLSServerConfig `json:"tls" yaml:"tls"`

	Probe *HTTPProbeConfig `json:"probe" yaml:"probe"`
}

type HTTPProbeConfig struct {
	// (optional) Hostname (expected Host header value) of the probe endpoint
	Host string `json:"host" yaml:"host"`
	// Path of the probe endpoint (e.g. "/probe")
	Path string `json:"path" yaml:"path"`
}
