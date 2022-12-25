package config

type PeeringConfig struct {
	Accept  *PeeringAcceptConfig    `json:"accept" yaml:"accept"`
	Connect []*PeeringConnectConfig `json:"connect" yaml:"connect"`
}

type PeeringAcceptConfig struct {
	// (optional) "host:port" to listen (e.g. "localhost:50051", ":50051")
	// Use random available port by default.
	Listen string `json:"listen" yaml:"listen"`

	TLS *TLSServerConfig `json:"tls" yaml:"tls"`
}

const PeeringConnectionDefaultCount = 3

type PeeringConnectConfig struct {
	/** e.g. "localhost:50051" */
	Address string `json:"address" yaml:"address"`

	TLS *TLSClientConfig `json:"tls" yaml:"tls"`

	Connections                int `json:"connections" yaml:"connections"`
	ConnectionRetryIntervalSec int `json:"connection_retry_interval_sec" yaml:"connection_retry_interval_sec"`
}
