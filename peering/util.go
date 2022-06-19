package peering

import (
	"context"
	"sync/atomic"

	"github.com/saiya/mesh_for_home_server/logger"
	"google.golang.org/grpc/peer"
)

// neverFail panics program only if non-nil error given
func neverFail(err error) {
	if err == nil {
		return
	}

	panic(err)
}

var connectionIDGen = uint64(0)

func withConnectionLogAttributes(ctx context.Context) context.Context {
	attributes := make([]interface{}, 0, 2*6)
	attributes = append(attributes, "connection-id", atomic.AddUint64(&connectionIDGen, 1))

	p, ok := peer.FromContext(ctx)
	if ok {
		attributes = append(
			attributes,
			"remote-addr", p.Addr.String(),
			"auth-type", p.AuthInfo.AuthType(),
		)
	}

	return logger.Wrap(ctx, attributes...)
}
