package config

type IngressConfigs struct {
	HTTP []HTTPIngressConfig `json:"http" yaml:"http"`
}
