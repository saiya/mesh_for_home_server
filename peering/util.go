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

func debugLogIfErr(ctx context.Context, err error) {
	if err == nil {
		return
	}
	logger.GetFrom(ctx).Debug(err)
}

var connectionIDGen = uint64(0)

type connectionIDKey struct{}

func withConnectionLogAttributes(ctx context.Context) context.Context {
	connectionID := ctx.Value(connectionIDKey{})
	if connectionID == nil {
		connectionID = atomic.AddUint64(&connectionIDGen, 1)
		ctx = context.WithValue(ctx, connectionIDKey{}, connectionID)
	}

	attributes := make([]interface{}, 0, 2*6)
	attributes = append(attributes, "connection-id", connectionID)

	p, ok := peer.FromContext(ctx)
	if ok {
		if p.Addr != nil {
			attributes = append(attributes, "remote-addr", p.Addr.String())
		}
		if p.AuthInfo != nil {
			attributes = append(attributes, "auth-type", p.AuthInfo.AuthType())
		}
	}

	return logger.Wrap(ctx, attributes...)
}
