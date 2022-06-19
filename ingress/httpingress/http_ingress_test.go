package httpingress

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/tlshelper"
	"github.com/saiya/mesh_for_home_server/tlshelper/tlstesting"
	"github.com/stretchr/testify/assert"
)

func TestHTTPIngress(t *testing.T) {
	ingress, err := NewHTTPIngress(
		&config.HTTPIngressConfig{
			Listen: "localhost:0",
			Probe: &config.HTTPProbeConfig{
				Host: "example.com",
				Path: "/probe",
			},
		},
		handlerMustNotBeCalled(t),
	)
	assert.NoError(t, err)
	defer ingress.Close(context.Background())

	testIngress(t, ingress.(*httpIngress), http.DefaultClient, nil)
}

func TestHTTPSIngress(t *testing.T) {
	ingress, err := NewHTTPIngress(
		&config.HTTPIngressConfig{
			Listen: "localhost:0",
			TLS:    tlstesting.GenerateServerCert("localhost"),
			Probe: &config.HTTPProbeConfig{
				Host: "example.com",
				Path: "/probe",
			},
		},
		handlerMustNotBeCalled(t),
	)
	assert.NoError(t, err)
	defer ingress.Close(context.Background())

	testIngress(
		t, ingress.(*httpIngress),
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: tlstesting.RootCAs(),
				},
			},
		},
		nil,
	)
}

func TestHTTPSIngressWithClientCert(t *testing.T) {
	testHTTPSIngressWithClientCert(t, false)
}

func TestHTTPSIngressWithClientCert_NoClientCert(t *testing.T) {
	testHTTPSIngressWithClientCert(t, true)
}

func testHTTPSIngressWithClientCert(t *testing.T, probeWithoutClientCert bool) {
	tlsConfig := tlstesting.GenerateServerCert("localhost")
	clientTLSConfig := tlstesting.EnableClientCert(tlsConfig)
	ingress, err := NewHTTPIngress(
		&config.HTTPIngressConfig{
			Listen: "localhost:0",
			TLS:    tlsConfig,
			Probe: &config.HTTPProbeConfig{
				Host: "example.com",
				Path: "/probe",
			},
		},
		handlerMustNotBeCalled(t),
	)
	assert.NoError(t, err)
	defer ingress.Close(context.Background())

	clientTLSRawConfig, err := tlshelper.BuildTLSClientConfig(clientTLSConfig)
	assert.NoError(t, err)
	if probeWithoutClientCert {
		clientTLSRawConfig.Certificates = []tls.Certificate{}
	}
	testIngress(
		t, ingress.(*httpIngress),
		&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: clientTLSRawConfig,
			},
		},
		func(err error) {
			if probeWithoutClientCert {
				assert.ErrorContains(t, err, "remote error: tls: bad certificate")
			} else {
				assert.NoError(t, err)
			}
		},
	)
}

func testIngress(t *testing.T, ingress *httpIngress, client *http.Client, probeError func(error)) {
	protocol := "http"
	tlsDescription := "none"
	if ingress.config.TLS != nil {
		protocol = "https"
		tlsDescription = "enabled"
		if ingress.config.TLS.ClientCertCAFile != "" {
			tlsDescription = "require_client_cert"
		}
	}

	port := ingress.port
	assert.Equal(t, fmt.Sprintf("HTTPIngress{ port: %d, TLS: %s }", port, tlsDescription), ingress.String())

	if ingress.config.Probe != nil {
		probeURL := fmt.Sprintf("%s://localhost:%d%s", protocol, port, ingress.config.Probe.Path)
		logger.Get().Infow("Testing probe endpoint", "url", probeURL)

		req, err := http.NewRequest("GET", probeURL, nil)
		assert.NoError(t, err)
		req.Host = ingress.config.Probe.Host

		res, err := client.Do(req)
		if probeError != nil {
			probeError(err)
		} else {
			assert.NoError(t, err)
		}
		if err == nil {
			assert.Equal(t, 200, res.StatusCode)
		}
	}
}

func TestHTTPListenFailure(t *testing.T) {
	_, err := NewHTTPIngress(
		&config.HTTPIngressConfig{
			Listen: "localhost:999999",
		},
		handlerMustNotBeCalled(t),
	)
	assert.ErrorContains(t, err, "failed to listen \"localhost:999999\"")
}

func handlerMustNotBeCalled(t *testing.T) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		t.Error("No handler call expected")
		return nil
	}
}
