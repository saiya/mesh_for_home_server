package ingress

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
)

var serverCertPublicKey *rsa.PublicKey
var serverCertPrivateKey *rsa.PrivateKey
var serverCertPrivateKeyFile string
var serverCertSerialNumber = big.NewInt(1)

var rootCA *x509.Certificate
var rootCAs = x509.NewCertPool() // For tls.Config.RootCAs

func init() {
	var err error

	serverCertPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	neverFail(err)
	serverCertPublicKey = &serverCertPrivateKey.PublicKey

	privBytes, err := x509.MarshalPKCS8PrivateKey(serverCertPrivateKey)
	neverFail(err)
	keyOut, err := ioutil.TempFile("", "key*.pem")
	neverFail(err)
	serverCertPrivateKeyFile = keyOut.Name()
	neverFail(pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}))
	neverFail(keyOut.Close())

	rootCA = generateRootCACert()
	rootCAs.AddCert(rootCA)
}

// ref: https://go.dev/src/crypto/tls/generate_cert.go
func generateServerCert(hostname string) *config.TLSIngressConfig {
	c := config.TLSIngressConfig{}

	serverCertSerialNumber.Add(serverCertSerialNumber, big.NewInt(1))
	template := x509.Certificate{
		SerialNumber:          serverCertSerialNumber,
		BasicConstraintsValid: true,

		Subject: pkix.Name{
			Organization: []string{"Dummy Value"},
		},
		DNSNames: []string{hostname},

		NotBefore: time.Now().Add(time.Second * -30),
		NotAfter:  time.Now().Add(time.Hour),

		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, rootCA, serverCertPublicKey, serverCertPrivateKey)
	neverFail(err)
	certOut, err := ioutil.TempFile("", "cert*.pem")
	neverFail(err)
	c.CertFile = certOut.Name()
	neverFail(pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))
	neverFail(certOut.Close())

	c.KeyFile = serverCertPrivateKeyFile
	return &c
}

// generateClientCert retruns (ClientCertCAFile, ClientCert) pair
func generateClientCert() (string, *tls.Certificate) {
	rootCA = generateRootCACert()

	rootCAOut, err := ioutil.TempFile("", "client-cert-root*.pem")
	neverFail(err)
	rootCACertFile := rootCAOut.Name()
	neverFail(pem.Encode(rootCAOut, &pem.Block{Type: "CERTIFICATE", Bytes: rootCA.Raw}))
	neverFail(rootCAOut.Close())

	clientCert, err := NewClientCert(
		rand.Reader, time.Now(),
		ClientCertGenerateParameter{
			pkcs12password: "test",

			ttl:          time.Hour,
			serialNumber: big.NewInt(1),
			commonName:   "Test client cert",

			rootCACertFile: rootCACertFile,
			rootCAKeyFile:  serverCertPrivateKeyFile,
		},
	)
	neverFail(err)

	return rootCACertFile, clientCert.TLS
}

func generateRootCACert() *x509.Certificate {
	rootCATemplate := x509.Certificate{
		SerialNumber:          big.NewInt(0),
		BasicConstraintsValid: true,

		IsCA: true,
		Subject: pkix.Name{
			Organization: []string{"Dummy Value"},
		},

		NotBefore: time.Now().Add(time.Second * -30),
		NotAfter:  time.Now().Add(time.Hour),

		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &rootCATemplate, &rootCATemplate, serverCertPublicKey, serverCertPrivateKey)
	neverFail(err)
	rootCA, err = x509.ParseCertificate(derBytes)
	neverFail(err)
	return rootCA
}
