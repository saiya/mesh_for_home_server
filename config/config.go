package config

type NodeID string

type ServerConfig struct {
	NodeID NodeID `json:"id"`

	Perring *PeeringConfig `json:"peering"`
}

type PeeringConfig struct {
	Accept  *PerringAcceptConfig   `json:"accept"`
	Connect []PeeringConnectConfig `json:"connect"`
}

type PerringAcceptConfig struct {
	PublicKeyJwk []string `json:"public_key_jwk"`
}

type PeeringConnectConfig struct {
	PrivateKeyJwk string `json:"private_key_jwk"`
}
