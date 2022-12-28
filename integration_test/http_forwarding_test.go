package integration_test_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/server"
	util "github.com/saiya/mesh_for_home_server/testutil"
	"github.com/stretchr/testify/assert"
)

func TestRequests(t *testing.T) {
	httpServer, destPort := util.NewHTTPStubServer(t, util.DEFAULT_HTTP_STUBS...)
	defer httpServer.Close()

	test(t, destPort, func(httpClient *http.Client, srcPort int) {
		for _, c := range util.DEFAULT_HTTP_TEST_CASES {
			t.Run(c.String(), func(t *testing.T) { c.Do(t, httpClient, srcPort) })
		}
	})
}

func test(t *testing.T, destPort int, f func(httpClient *http.Client, srcPort int)) {
	srv1IngressPort := 8089
	srv2PeerPort := 8088

	srv1, err := server.StartServer(
		&config.ServerConfig{
			Hostname: "ingress",
			Perring: &config.PeeringConfig{
				Connect: []*config.PeeringConnectConfig{
					{Address: fmt.Sprintf("localhost:%d", srv2PeerPort), Connections: 1},
				},
			},
			Ingress: &config.IngressConfigs{
				HTTP: []config.HTTPIngressConfig{
					{
						Listen:          fmt.Sprintf("0.0.0.0:%d", srv1IngressPort),
						RequestTimeout:  &config.HTTPTimeout{Header: "1s", Body: "10s"},
						ResponseTimeout: &config.HTTPTimeout{Header: "1s", Body: "10s"},
					},
				},
			},
		},
	)
	if !assert.NoError(t, err) {
		return
	}
	defer srv1.Close(context.Background())

	srv2, err := server.StartServer(
		&config.ServerConfig{
			Hostname: "egress",
			Perring: &config.PeeringConfig{
				Accept: &config.PeeringAcceptConfig{
					Listen: fmt.Sprintf("localhost:%d", srv2PeerPort),
				},
			},
			Egress: &config.EgressConfigs{
				HTTP: []config.HTTPEgressConfig{
					{
						Host: "*", Server: fmt.Sprintf("http://localhost:%d", destPort),
						ConnectTimeout:  "1s",
						ResponseTimeout: &config.HTTPTimeout{Header: "1s", Body: "10s"},
					},
				},
			},
		},
	)
	if !assert.NoError(t, err) {
		return
	}
	defer srv2.Close(context.Background())

	time.Sleep(500 * time.Millisecond) // FIXME: (hack) Await peering & advertisement
	f(&http.Client{}, srv1IngressPort)
}
