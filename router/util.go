package router

import (
	"context"

	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

func withRoutingMessageLogContext(ctx context.Context, routing *generated.RoutingMessage) context.Context {
	ctx = logger.Wrap(
		ctx,
		"dest", routing.DestNodeId, "via", routing.ViaNodeId,
	)
	if routing.GetHttp() != nil {
		ctx = logger.Wrap(ctx, "http_hostname_pattern", routing.GetHttp().GetHostnamePattern())
	}
	return ctx
}
