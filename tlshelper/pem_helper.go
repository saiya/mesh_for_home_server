package tlshelper

import (
	"crypto"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"strings"
)

// ParseCertificatePemFile expects PEM(DER(X509(cert))) format
func ParseCertificatePemFile(filepath string) (*x509.Certificate, error) {
	pemBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("unable to load cert file: %w", err)
	}
	pemBlock, _ := pem.Decode(pemBytes)
	if pemBlock.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("unable to load cert file: PEM file must have only CERTIFICATE block")
	}
	cert, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse cert file: %w", err)
	}
	return cert, nil
}

// ParsePrivateKeyPemFile expects PEM(DER(PKCS8(key))) format
func ParsePrivateKeyPemFile(filepath string) (crypto.PrivateKey, error) {
	pemBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("Unable to load client certificate root CA key file: %w", err)
	}
	pemBlock, _ := pem.Decode(pemBytes)
	if !(pemBlock.Type == "PRIVATE KEY" || strings.HasSuffix(pemBlock.Type, " PRIVATE KEY")) {
		return nil, fmt.Errorf("Unable to load client certificate root CA key file: PEM file must have only PRIVATE KEY block")
	}
	privKeyObj, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client certificate root CA key file: %w", err)
	}
	if privKey, ok := privKeyObj.(crypto.PrivateKey); ok {
		return privKey, nil
	}
	return nil, fmt.Errorf("expected crypto.PrivateKey but loaded different object")
}
