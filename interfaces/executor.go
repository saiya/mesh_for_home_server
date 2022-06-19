package interfaces

import (
	"context"
	"net/http"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
)

type Executor interface {
	Close(context.Context) error
}

type HTTPExecutor interface {
	Executor
	http.RoundTripper
}

type PingExecutor interface {
	Executor

	Ping(ctx context.Context, dest config.NodeID, timeout time.Duration) error
}
