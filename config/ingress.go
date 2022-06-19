package config

type IngressConfig struct {
	HTTP *HTTPIngressConfig `json:"http"`
}

type HTTPIngressConfig struct {
	// (optional) "host:port" to listen (e.g. "localhost:8080", ":8080")
	// Use random available port by default.
	Listen string `json:"listen"`

	TLS *TLSIngressConfig `json:"tls"`

	Probe *HTTPProbeConfig `json:"probe"`
}

type TLSIngressConfig struct {
	// (optional) "TLS1.2" or "TLS1.3"
	MinVersion string `json:"min_version"`
	// (optional) "TLS1.2" or "TLS1.3"
	MaxVersion string `json:"max_version"`

	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`

	// (optional) If present, require client certificate
	ClientCertCAFile string `json:"client_cert_ca_file"`
}

type HTTPProbeConfig struct {
	// (optional) Hostname (expected Host header value) of the probe endpoint
	Host string
	// Path of the probe endpoint (e.g. "/probe")
	Path string
}
