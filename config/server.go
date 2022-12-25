package config

type ServerConfig struct {
	Perring *PeeringConfig  `json:"peering" yaml:"peering"`
	Ingress *IngressConfigs `json:"ingress" yaml:"ingress"`
	Egress  *EgressConfigs  `json:"egress" yaml:"egress"`
}
