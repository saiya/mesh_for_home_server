package config

import "time"

type EgressConfigs struct {
	AdvertiseIntervalSec int `json:"advertise-interval-sec" yaml:"advertise-interval-sec"`

	HTTP []HTTPEgressConfig `json:"http" yaml:"http"`
}

const AdvertiseIntervalDefault = 15 * time.Second
const AdvertiseTTLMargin = 10 * time.Second
