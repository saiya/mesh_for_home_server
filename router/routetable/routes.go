package routetable

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
)

type Routes struct {
	lock   sync.RWMutex
	routes map[config.NodeID]Route
}

func NewRoutes() *Routes {
	return &Routes{
		routes: make(map[config.NodeID]Route),
	}
}

func (rs *Routes) Save(r Route) {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.routes[r.Dest] = r
}

// RoundRobin over valid routes, or returns ErrNoRoute
func (rs *Routes) RoundRobin(ctx context.Context, now time.Time) (config.NodeID, error) {
	triggerGC := false

	rs.lock.RLock()
	defer func() {
		rs.lock.RUnlock()
		if triggerGC {
			rs.GC(now)
		}
	}()

	candidates := make([]Route, 0, len(rs.routes))
	for _, r := range rs.routes {
		if r.IsValid(now) {
			candidates = append(candidates, r)
		} else {
			triggerGC = true
			logger.GetFrom(ctx).Debugw("Stale/invalid route found", "dest", r.Dest, "expireAt", r.ExpireAt)
		}
	}

	if len(candidates) == 0 {
		return "", interfaces.ErrNoRoute
	}
	return candidates[rand.Intn(len(candidates))].Dest, nil
}

func (rs *Routes) GC(now time.Time) {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	for dest, r := range rs.routes {
		if !r.IsValid(now) {
			delete(rs.routes, dest)
		}
	}
}
