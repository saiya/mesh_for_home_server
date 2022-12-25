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

type httpRT struct {
	lock                sync.RWMutex
	routesByHostPattern map[string]*httpRoute
}

func NewHTTPRT() *httpRT {
	return &httpRT{
		routesByHostPattern: map[string]*httpRoute{},
	}
}

func (rt *httpRT) Close(ctx context.Context) error {
	return nil
}

func (rt *httpRT) Update(ctx context.Context, node config.NodeID, expireAt time.Time, ad *generated.HttpAdvertisement) {
	rt.lock.Lock()
	defer rt.lock.Unlock()

	for _, hostPattern := range ad.HostnameMatchers {
		route := rt.routesByHostPattern[hostPattern]
		if route == nil {
			route = newHttpRoute(hostPattern)
			rt.routesByHostPattern[hostPattern] = route
		}
		route.routes.Save(Route{ExpireAt: expireAt, Dest: node})
	}
}

func (rt *httpRT) Route(ctx context.Context, now time.Time, request *generated.HttpRequestStart) (config.NodeID, error) {
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

func newHttpRoute(hostPattern string) *httpRoute {
	return &httpRoute{
		hostPattern:        hostPattern,
		hostPatternMatcher: dnshelper.HostnameMatcher(hostPattern),
		routes:             NewRoutes(),
	}
}
