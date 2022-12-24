package config

type NodeID string

type ServerConfig struct {
	NodeID NodeID `json:"id"`

	Perring *PeeringConfig  `json:"peering"`
	Ingress *IngressConfigs `json:"ingress"`
	Egress  *EgressConfigs  `json:"egress"`
}
