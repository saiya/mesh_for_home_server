package integration_test_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/egress/handler/httphandler"
	"github.com/saiya/mesh_for_home_server/ingress/forwarder"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/router"
	util "github.com/saiya/mesh_for_home_server/testutil"
)

func TestRequests(t *testing.T) {
	httpServer, port := util.NewHTTPStubServer(t, util.DEFAULT_HTTP_STUBS...)
	defer httpServer.Close()

	test(t, port, func(httpClient *http.Client) {
		for _, c := range util.DEFAULT_HTTP_TEST_CASES {
			t.Run(c.String(), func(t *testing.T) { c.Do(t, httpClient, 80) })
		}
	})
}

func test(t *testing.T, port int, f func(httpClient *http.Client)) {
	ctx := util.Context(t)

	rt := router.NewRouter("test")
	defer rt.Close(ctx)

	ingress, httpClient := startIngress(rt)
	defer ingress.Close(ctx)

	egress := httphandler.NewHttpHandler(rt)
	defer egress.Close(ctx)
	egress.AddEgress(&config.HTTPEgressConfig{
		Host:   "*",
		Server: fmt.Sprintf("localhost:%d", port),
	})

	f(httpClient)
}

func startIngress(rt interfaces.Router) (interfaces.HTTPForwarder, *http.Client) {
	ingress := forwarder.NewHTTPForwarder(rt)
	return ingress, &http.Client{
		Transport: ingress,
		Timeout:   1 * time.Second,
	}
}
