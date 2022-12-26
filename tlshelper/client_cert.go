package tlshelper

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"os"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"software.sslmate.com/src/go-pkcs12"
)

type ClientCertGenerateParameter struct {
	PKCS12password string

	TTL          time.Duration
	SerialNumber *big.Int
	CommonName   string

	RootCACertFile string
	RootCAKeyFile  string
}

// ParsePKCS12File loads PKCS#12 encoded client credentials (certificate + private key)
func ParsePKCS12File(config config.MTLSCertLoadConfig) (*tls.Certificate, error) {
	pfxData, err := os.ReadFile(config.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed to read client credentials file (%s): %w", config.Path, err)
	}

	privKey, cert, err := pkcs12.Decode(pfxData, config.Password)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse PKCE#12 encoded credentials (%s): %w", config.Path, err)
	}

	return &tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privKey.(crypto.PrivateKey),
		Leaf:        cert,
	}, nil
}

type ClientCertGenerateResult struct {
	// PKCS#12 data contains both private key and client certificate
	PKCS12 []byte

	// TLS certificate object, contains both private key and client certificate
	// Can be used in combination with tls.Config and Go HTTP client.
	TLS *tls.Certificate
}

func NewClientCert(
	rng io.Reader, now time.Time,
	params ClientCertGenerateParameter,
) (ClientCertGenerateResult, error) {
	rootCACert, err := ParseCertificatePemFile(params.RootCACertFile)
	if err != nil {
		return ClientCertGenerateResult{}, fmt.Errorf("Unable to load client certificate root CA cert file: %w", err)
	}

	rootCAPrivateKey, err := ParsePrivateKeyPemFile(params.RootCAKeyFile)
	if err != nil {
		return ClientCertGenerateResult{}, fmt.Errorf("Unable to load client certificate root CA key file: %w", err)
	}

	privKey, err := ecdsa.GenerateKey(elliptic.P521(), rng)
	neverFail(err)
	pubKey := &privKey.PublicKey
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	neverFail(err)

	clientCertBytes, err := x509.CreateCertificate(
		rng,
		&x509.Certificate{
			SerialNumber:          params.SerialNumber,
			BasicConstraintsValid: true,

			Subject: pkix.Name{
				CommonName: params.CommonName,
			},

			NotBefore: now,
			NotAfter:  now.Add(params.TTL),

			// For RSA, also x509.KeyUsageKeyEncipherment needed
			KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment,
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageClientAuth,
				x509.ExtKeyUsageCodeSigning,
			},
		},
		rootCACert,
		pubKey,
		rootCAPrivateKey,
	)
	neverFail(err)

	clientCert, err := x509.ParseCertificate(clientCertBytes)
	neverFail(err)

	pfxBytes, err := pkcs12.Encode(rng, privKey, clientCert, []*x509.Certificate{}, params.PKCS12password)
	neverFail(err)

	tlsKeyPair, err := tls.X509KeyPair(
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCertBytes}),
		pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privKeyBytes}),
	)
	neverFail(err)

	return ClientCertGenerateResult{
		PKCS12: pfxBytes,
		TLS:    &tlsKeyPair,
	}, nil
}
