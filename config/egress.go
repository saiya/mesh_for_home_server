package config

import "time"

type EgressConfigs struct {
	AdvertiseIntervalSec int `json:"advertise-interval-sec" yaml:"advertise-interval-sec"`

	HTTP []HTTPEgressConfig `json:"http" yaml:"http"`
}

const AdvertiseIntervalDefault = 15 * time.Second
const AdvertiseTtlMargin = 10 * time.Second

type HTTPEgressConfig struct {
	// (optional) This egress receive HTTP traffic that match with this hostname.
	// By default or for empty string, match with any host.
	//
	// Can use "*" in the 1st part of hostname (e.g. "*.example.com" match with "sub.example.com").
	// "*" match with only 1 level of hostname element (e.g. "*.example.com" not match with "a.b.example.com")
	//
	// If multiple egress matches with request, longest match preffered.
	Host string `json:"host" yaml:"host"`

	// Upstream serer, "hostname:port" format (e.g. "localhost:8080")
	Server string `json:"server" yaml:"server"`
}
