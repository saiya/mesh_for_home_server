package tlstesting

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/tlshelper"
)

// For performance, reuse same key pair for all test certs
var serverCertPublicKey *rsa.PublicKey
var serverCertPrivateKey *rsa.PrivateKey
var serverCertPrivateKeyFile string

var rootCA *tlshelper.RootCAGenerateResult

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

	tempdir, err := ioutil.TempDir("", "test-server-cert")
	neverFail(err)
	rootCA, err = tlshelper.NewSelfSignedRootCA(
		&tlshelper.RootCAGenerateRequest{
			RNG:          rand.Reader,
			NotBefore:    time.Now().Add(-time.Second * 3),
			NotAfter:     time.Now().Add(time.Hour),
			FilePathBase: tempdir + "/rootCA",
			CommonName:   "Test client cert Root CA",
		},
	)
	neverFail(err)
}

func RootCAs() *x509.CertPool {
	return rootCA.CACertAsPool()
}

var serverCertSerialNumber = big.NewInt(1)

// ref: https://go.dev/src/crypto/tls/generate_cert.go
func GenerateServerCert(hostname string) *config.TLSServerConfig {
	c := config.TLSServerConfig{}

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

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, rootCA.CACert, serverCertPublicKey, rootCA.PrivateKey)
	neverFail(err)
	certOut, err := ioutil.TempFile("", "cert*.pem")
	neverFail(err)
	c.CertFile = certOut.Name()
	neverFail(pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))
	neverFail(certOut.Close())

	c.KeyFile = serverCertPrivateKeyFile
	return &c
}

func EnableClientCert(serverConfig *config.TLSServerConfig) *config.TLSClientConfig {
	clientCARoot, clientCert := generateClientCert()
	clientCertFile, err := ioutil.TempFile("", "client-cert*.p12")
	neverFail(err)
	_, err = clientCertFile.Write(clientCert.PKCS12)
	neverFail(err)
	neverFail(clientCertFile.Close())

	serverConfig.ClientCertCAFile = clientCARoot.CACertFile
	return &config.TLSClientConfig{
		RootCAFiles: []string{rootCA.CACertFile},
		ClientCerts: []config.MTLSCertLoadConfig{
			{
				Path:     clientCertFile.Name(),
				Password: "test",
			},
		},
	}
}

func generateClientCert() (*tlshelper.RootCAGenerateResult, tlshelper.ClientCertGenerateResult) {
	tempdir, err := ioutil.TempDir("", "test-client-cert")
	neverFail(err)

	clientCertRootCA, err := tlshelper.NewSelfSignedRootCA(
		&tlshelper.RootCAGenerateRequest{
			RNG:          rand.Reader,
			NotBefore:    time.Now().Add(-time.Second * 3),
			NotAfter:     time.Now().Add(time.Hour),
			FilePathBase: tempdir + "/rootCA",
			CommonName:   "Test client cert Root CA",
		},
	)
	neverFail(err)

	clientCert, err := tlshelper.NewClientCert(
		rand.Reader, time.Now(),
		tlshelper.ClientCertGenerateParameter{
			PKCS12password: "test",

			TTL:          time.Hour,
			SerialNumber: big.NewInt(1),
			CommonName:   "Test client cert",

			RootCACertFile: clientCertRootCA.CACertFile,
			RootCAKeyFile:  clientCertRootCA.PrivateKeyFile,
		},
	)
	neverFail(err)

	return clientCertRootCA, clientCert
}
