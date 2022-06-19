package executor

import (
	"context"
	"fmt"
	"net/http"

	"github.com/saiya/mesh_for_home_server/interfaces"
)

type httpExecutor struct{}

func NewHTTPExecutor(router interfaces.Router) interfaces.HTTPExecutor {
	// FIXME: Implement
	return &httpExecutor{}
}

func (exec *httpExecutor) Close(ctx context.Context) error {
	// FIXME: Implement
	return nil
}

func (exec *httpExecutor) RoundTrip(req *http.Request) (*http.Response, error) {
	// FIXME: Implement
	return nil, fmt.Errorf("Not implemented yet")
}
