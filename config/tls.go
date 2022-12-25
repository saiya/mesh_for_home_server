package config

type TLSServerConfig struct {
	// (optional) "TLS1.2" or "TLS1.3"
	MinVersion string `json:"min_version" yaml:"min_version"`
	// (optional) "TLS1.2" or "TLS1.3"
	MaxVersion string `json:"max_version" yaml:"max_version"`

	CertFile string `json:"cert_file" yaml:"cert_file"`
	KeyFile  string `json:"key_file" yaml:"key_file"`

	// (optional) If present, require client certificate (mTLS)
	ClientCertCAFile string `json:"client_cert_ca_file" yaml:"client_cert_ca_file"`
}

type TLSClientConfig struct {
	// (optional) PKCS#12 encoded client credentials for mTLS
	ClientCerts []MTLSCertLoadConfig `json:"client_certs" yaml:"client_certs"`

	// (optional) If present, use those certs for Root CA instead of environment provided root CA list.
	RootCAFiles []string `json:"root_ca_files" yaml:"root_ca_files"`
}

type MTLSCertLoadConfig struct {
	// Path of PKCS#12 encoded credential file
	Path string `json:"path" yaml:"path"`
	// Password of PKCS#12 archive
	Password string `json:"password" yaml:"password"`
}
