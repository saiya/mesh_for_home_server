package forwarder

import (
	"context"
	"fmt"
	"net/http"

	"github.com/saiya/mesh_for_home_server/interfaces"
)

type httpForwarder struct{}

func NewHTTPForwarder(router interfaces.Router) interfaces.HTTPForwarder {
	// FIXME: Implement
	return &httpForwarder{}
}

func (exec *httpForwarder) Close(ctx context.Context) error {
	// FIXME: Implement
	return nil
}

func (exec *httpForwarder) RoundTrip(req *http.Request) (*http.Response, error) {
	// FIXME: Implement
	return nil, fmt.Errorf("Not implemented yet")
}
