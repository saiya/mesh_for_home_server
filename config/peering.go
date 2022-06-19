package config

type PeeringConfig struct {
	Accept  *PeeringAcceptConfig   `json:"accept"`
	Connect []PeeringConnectConfig `json:"connect"`
}

type PeeringAcceptConfig struct {
	// (optional) "host:port" to listen (e.g. "localhost:50051", ":50051")
	// Use random available port by default.
	Listen string `json:"listen"`

	TLS *TLSServerConfig `json:"tls"`
}

type PeeringConnectConfig struct {
	/** e.g. "localhost:50051" */
	Address string

	TLS *TLSClientConfig `json:"tls"`

	ConnectionRetryIntervalSec int `json:"connection_retry_interval_sec"`
}
