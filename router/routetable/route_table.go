package routetable

import (
	"context"
	"fmt"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
)

type RouteTable interface {
	Close(ctx context.Context) error

	Update(ctx context.Context, node config.NodeID, ad interfaces.Advertisement)
	Route(ctx context.Context, request interfaces.Message) (config.NodeID, error)
}

type routeTable struct {
	http *httpRT
}

func NewRouteTable() RouteTable {
	return &routeTable{
		http: NewHTTPRT(),
	}
}

func (rt *routeTable) Close(ctx context.Context) error {
	return rt.http.Close(ctx)
}

func (rt *routeTable) Update(ctx context.Context, node config.NodeID, ad interfaces.Advertisement) {
	expireAt := time.Unix(ad.ExpireAt, 0)
	rt.http.Update(ctx, node, expireAt, ad.Http)
}

func (rt *routeTable) Route(ctx context.Context, request interfaces.Message) (config.NodeID, error) {
	if http := request.GetHttpRequestStart(); http != nil {
		dest, err := rt.http.Route(ctx, time.Now(), http)
		logger.GetFrom(ctx).Debugw("HTTP routing done", "dest", dest, "success", err != nil)
		return dest, err
	}
	return "", fmt.Errorf("Cannot route given type of message: %v", request)
}