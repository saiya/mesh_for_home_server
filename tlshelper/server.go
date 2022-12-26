package tlshelper

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/saiya/mesh_for_home_server/config"
)

func BuildTLSServerConfig(cfg *config.TLSServerConfig) (*tls.Config, error) {
	c := tls.Config{
		MinVersion: parseTLSVersion(cfg.MinVersion, tls.VersionTLS12),
		MaxVersion: parseTLSVersion(cfg.MinVersion, tls.VersionTLS13),
	}
	if err := setupServerCert(&c, cfg); err != nil {
		return nil, err
	}
	if err := setupTLSClientAuth(&c, cfg); err != nil {
		return nil, err
	}
	return &c, nil
}

func setupServerCert(c *tls.Config, cfg *config.TLSServerConfig) error {
	var err error
	c.Certificates = make([]tls.Certificate, 1)
	c.Certificates[0], err = tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		return fmt.Errorf("failed to load server certificate: %w", err)
	}
	return nil
}

func setupTLSClientAuth(c *tls.Config, cfg *config.TLSServerConfig) error {
	if cfg.ClientCertCAFile != "" {
		caCertFile, err := os.ReadFile(cfg.ClientCertCAFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS client cert CA file (%v): %w", cfg.ClientCertCAFile, err)
		}

		c.ClientAuth = tls.RequireAndVerifyClientCert
		c.ClientCAs = x509.NewCertPool()
		c.ClientCAs.AppendCertsFromPEM(caCertFile)
	}
	return nil
}
