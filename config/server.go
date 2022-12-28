package config

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Hostname string `json:"hostname" yaml:"hostname"` // = prefix of NodeID

	Perring *PeeringConfig  `json:"peering" yaml:"peering"`
	Ingress *IngressConfigs `json:"ingress" yaml:"ingress"`
	Egress  *EgressConfigs  `json:"egress" yaml:"egress"`
}

func ParseConfig(in io.Reader) (*ServerConfig, error) {
	src, err := io.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &ServerConfig{}
	err = yaml.Unmarshal(src, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file YAML: %w", err)
	}
	return config, nil
}
