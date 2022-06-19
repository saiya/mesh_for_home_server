package httpegress

import (
	"context"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
)

type httpEgress struct {
}

func NewHTTPEgress(c *config.HTTPEgressConfig) (interfaces.Egress, error) {
	return &httpEgress{}, nil
}

func (e *httpEgress) String() string {
	// FIXME: Implement
	panic("Not implemented yet")
}

func (e *httpEgress) Close(ctx context.Context) error {
	// FIXME: Implement
	panic("Not implemented yet")
}
