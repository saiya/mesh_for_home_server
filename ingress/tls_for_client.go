package ingress

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"strings"
	"time"

	"software.sslmate.com/src/go-pkcs12"
)

type ClientCertGenerateParameter struct {
	pkcs12password string

	ttl          time.Duration
	serialNumber *big.Int
	commonName   string

	rootCACertFile string
	rootCAKeyFile  string
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
	rootCACertBytes, err := ioutil.ReadFile(params.rootCACertFile)
	if err != nil {
		return ClientCertGenerateResult{}, fmt.Errorf("Unable to load client certificate root CA cert file: %w", err)
	}
	rootCACertPEMBlock, _ := pem.Decode(rootCACertBytes)
	if rootCACertPEMBlock.Type != "CERTIFICATE" {
		return ClientCertGenerateResult{}, fmt.Errorf("Unable to load client certificate root CA cert file: PEM file must have only CERTIFICATE block")
	}
	rootCACert, err := x509.ParseCertificate(rootCACertPEMBlock.Bytes)
	if err != nil {
		return ClientCertGenerateResult{}, fmt.Errorf("Unable to parse client certificate root CA cert file: %w", err)
	}

	rootCAKeyBytes, err := ioutil.ReadFile(params.rootCAKeyFile)
	if err != nil {
		return ClientCertGenerateResult{}, fmt.Errorf("Unable to load client certificate root CA key file: %w", err)
	}
	rootCAKeyPEMBlock, _ := pem.Decode(rootCAKeyBytes)
	if !(rootCAKeyPEMBlock.Type == "PRIVATE KEY" || strings.HasSuffix(rootCACertPEMBlock.Type, " PRIVATE KEY")) {
		return ClientCertGenerateResult{}, fmt.Errorf("Unable to load client certificate root CA key file: PEM file must have only PRIVATE KEY block")
	}
	rootCAPrivateKey, err := x509.ParsePKCS8PrivateKey(rootCAKeyPEMBlock.Bytes)
	if err != nil {
		return ClientCertGenerateResult{}, fmt.Errorf("Unable to parse client certificate root CA key file: %w", err)
	}

	privKey, err := ecdsa.GenerateKey(elliptic.P521(), rng)
	neverFail(err)
	pubKey := &privKey.PublicKey
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	neverFail(err)

	clientCertBytes, err := x509.CreateCertificate(
		rng,
		&x509.Certificate{
			SerialNumber:          params.serialNumber,
			BasicConstraintsValid: true,

			Subject: pkix.Name{
				CommonName: params.commonName,
			},

			NotBefore: now,
			NotAfter:  now.Add(params.ttl),

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

	pfxBytes, err := pkcs12.Encode(rng, privKey, clientCert, []*x509.Certificate{rootCACert}, params.pkcs12password)
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
