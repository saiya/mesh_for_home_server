package httphandler

import (
	"context"
	"net"
	"net/http"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/dnshelper"
)

type httpEgress struct {
	httpHandler     *httpHandler
	responseTimeout *config.HTTPTimeout

	host        string
	hostMatcher func(string) int64

	server     string // e.g. "http://localhost:8080"
	httpClient *http.Client
}

func newHTTPEgress(httpHandler *httpHandler, c *config.HTTPEgressConfig) *httpEgress {
	dialer := &net.Dialer{}
	c.ConfigureDialer(dialer)

	httpTransport := &http.Transport{
		DialContext:        dialer.DialContext,
		DisableCompression: true,
	}
	c.ConfigureHTTPTransport(httpTransport)

	return &httpEgress{
		httpHandler:     httpHandler,
		responseTimeout: c.ResponseTimeout,

		host:        c.Host,
		hostMatcher: dnshelper.HostnameMatcher(c.Host),

		server: c.Server,
		httpClient: &http.Client{
			Transport: httpTransport,
			Jar:       nil, // No automatic cookie handling
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // Don't follow any redirection
			},
		},
	}
}

func (h *httpEgress) Close(ctx context.Context) error {
	h.httpClient.CloseIdleConnections()
	return nil
}
