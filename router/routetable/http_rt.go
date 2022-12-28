package routetable

import (
	"context"
	"sync"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/dnshelper"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type HTTPRoutingTable struct {
	lock                sync.RWMutex
	routesByHostPattern map[string]*httpRoute
}

func NewHTTPRoutingTable() *HTTPRoutingTable {
	return &HTTPRoutingTable{
		routesByHostPattern: map[string]*httpRoute{},
	}
}

func (rt *HTTPRoutingTable) Close(ctx context.Context) error {
	return nil
}

func (rt *HTTPRoutingTable) Update(ctx context.Context, node config.NodeID, expireAt time.Time, ad *generated.HttpAdvertisement) {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	for _, hostPattern := range ad.HostnameMatchers {
		route := rt.routesByHostPattern[hostPattern]
		if route == nil {
			route = newHTTPRoute(hostPattern)
			rt.routesByHostPattern[hostPattern] = route
		}
		route.routes.Save(Route{ExpireAt: expireAt, Dest: node})
		logger.GetFrom(ctx).Debugw("HTTP route advertisement acknowledged", "host", hostPattern, "expire", expireAt, "dest", node)
	}
}

func (rt *HTTPRoutingTable) Route(ctx context.Context, now time.Time, request *generated.HttpRequestStart) (config.NodeID, error) {
	requestHostname := request.Hostname

	rt.lock.RLock()
	defer rt.lock.RUnlock()

	var route *httpRoute
	routePriority := int64(-1)
	for _, r := range rt.routesByHostPattern {
		priority := r.hostPatternMatcher(requestHostname)
		if priority < 0 { // Not match
			continue
		} else if route == nil || routePriority < priority {
			route = r
			routePriority = priority
		}
	}
	logger.GetFrom(ctx).Debugw("HTTP routing dispatched", "route", route, "routePriority", routePriority, "requestHostname", requestHostname)

	if route == nil {
		return "", interfaces.ErrNoRoute
	}
	return route.routes.RoundRobin(ctx, now)
}

type httpRoute struct {
	hostPattern        string
	hostPatternMatcher func(string) int64
	routes             *Routes
}

func (r *httpRoute) String() string {
	return r.hostPattern
}

func newHTTPRoute(hostPattern string) *httpRoute {
	return &httpRoute{
		hostPattern:        hostPattern,
		hostPatternMatcher: dnshelper.HostnameMatcher(hostPattern),
		routes:             NewRoutes(),
	}
}
