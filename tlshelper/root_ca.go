package tlshelper

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"time"
)

type RootCAGenerateRequest struct {
	RNG io.Reader

	// Filepath without suffix/extension.
	FilePathBase string

	NotBefore time.Time
	NotAfter  time.Time

	CommonName string
}

type RootCAGenerateResult struct {
	PublicKey crypto.PublicKey

	PrivateKeyFile string
	PrivateKey     crypto.PrivateKey

	CACertFile string
	CACert     *x509.Certificate
}

func (result *RootCAGenerateResult) CACertAsPool() *x509.CertPool {
	pool := x509.NewCertPool()
	pool.AddCert(result.CACert)
	return pool
}

func NewSelfSignedRootCA(req *RootCAGenerateRequest) (*RootCAGenerateResult, error) {
	result := RootCAGenerateResult{
		PrivateKeyFile: req.FilePathBase + ".key.pem",
		CACertFile:     req.FilePathBase + ".cert.pem",
	}

	err := generateRootCAKeys(req, &result)
	if err != nil {
		return nil, err
	}

	err = generateRootCASelfSignCert(req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func generateRootCAKeys(req *RootCAGenerateRequest, result *RootCAGenerateResult) error {
	var err error
	// result.PublicKey, result.PrivateKey, err := ed25519.GenerateKey(req.RNG)
	// neverFail(err)
	result.PrivateKey, err = ecdsa.GenerateKey(elliptic.P521(), req.RNG)
	neverFail(err)
	result.PublicKey = &result.PrivateKey.(*ecdsa.PrivateKey).PublicKey

	privKeyX509, err := x509.MarshalPKCS8PrivateKey(result.PrivateKey)
	neverFail(err)

	err = ioutil.WriteFile(
		result.PrivateKeyFile,
		pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privKeyX509}),
		0600,
	)
	if err != nil {
		return fmt.Errorf("Failed to write private key file: %w", err)
	}
	return nil
}

func generateRootCASelfSignCert(req *RootCAGenerateRequest, result *RootCAGenerateResult) error {
	rootCATemplate := x509.Certificate{
		SerialNumber:          big.NewInt(0),
		BasicConstraintsValid: true,

		IsCA: true,
		Subject: pkix.Name{
			CommonName: req.CommonName,
		},

		NotBefore: req.NotBefore,
		NotAfter:  req.NotAfter,

		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &rootCATemplate, &rootCATemplate, result.PublicKey, result.PrivateKey)
	neverFail(err)
	result.CACert, err = x509.ParseCertificate(derBytes)
	neverFail(err)

	err = ioutil.WriteFile(
		result.CACertFile,
		pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: result.CACert.Raw}),
		0644,
	)
	if err != nil {
		return fmt.Errorf("Failed to write CA cert file: %w", err)
	}
	return nil
}
