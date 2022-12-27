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

	ResponseTimeout *HTTPTimeout
}
