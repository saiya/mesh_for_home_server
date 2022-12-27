package interfaces

import (
	"context"
	"net/http"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
)

type Forwarder interface {
	Close(context.Context) error
}

type HTTPForwarder interface {
	Forwarder
	NewRoundTripper(cfg *config.HTTPIngressConfig) http.RoundTripper
}

type PingForwarder interface {
	Forwarder

	Ping(ctx context.Context, dest config.NodeID, timeout time.Duration) error
}
