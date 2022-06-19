package tlshelper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
)

func BuildTLSClientConfig(cfg *config.TLSClientConfig) (*tls.Config, error) {
	result := tls.Config{}

	for i, file := range cfg.ClientCerts {
		cert, err := ParsePKCS12File(file)
		if err != nil {
			return nil, fmt.Errorf("failed to load client cert [%d]: %w", i, err)
		}

		result.Certificates = append(result.Certificates, *cert)
	}

	for i, file := range cfg.RootCAFiles {
		cert, err := ParseCertificatePemFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to load Root CA cert [%d]: %w", i, err)
		}

		if result.RootCAs == nil {
			result.RootCAs = x509.NewCertPool()
		}
		result.RootCAs.AddCert(cert)
	}

	return &result, nil
}
