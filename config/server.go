package config

type NodeID string

type ServerConfig struct {
	NodeID NodeID `json:"id" yaml:"id"`

	Perring *PeeringConfig  `json:"peering" yaml:"peering"`
	Ingress *IngressConfigs `json:"ingress" yaml:"ingress"`
	Egress  *EgressConfigs  `json:"egress" yaml:"egress"`
}
