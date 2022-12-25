package router

import (
	"context"
	"math/rand"
	"sync"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"golang.org/x/sync/errgroup"
)

type outbounds struct {
	m     sync.Mutex
	table map[config.NodeID]*outbound
}

type outbound struct {
	m           sync.Mutex
	idGenerator uint64
	table       map[uint64]interfaces.RouterShink
}

func newOutbounds() *outbounds {
	return &outbounds{
		table: make(map[config.NodeID]*outbound),
	}
}

func newOutbound() *outbound {
	return &outbound{
		table: make(map[uint64]interfaces.RouterShink),
	}
}

func (u *outbounds) Register(dest config.NodeID, callback interfaces.RouterShink) interfaces.RouterUnregister {
	var unicast *outbound
	func() {
		u.m.Lock()
		defer u.m.Unlock()

		unicast = u.table[dest]
		if unicast == nil {
			unicast = newOutbound()
			u.table[dest] = unicast
		}
	}()
	return unicast.Register(callback)
}

func (u *outbound) Register(callback interfaces.RouterShink) interfaces.RouterUnregister {
	u.m.Lock()
	defer u.m.Unlock()

	u.idGenerator++
	id := u.idGenerator
	u.table[id] = callback
	return func() {
		u.m.Lock()
		defer u.m.Unlock()

		delete(u.table, id)
	}
}

func (u *outbounds) Broadcast(ctx context.Context, msg interfaces.Message) error {
	var eg errgroup.Group

	func() {
		u.m.Lock()
		defer u.m.Unlock()

		for _, unicast := range u.table {
			eg.Go(func() error {
				return unicast.Deliver(ctx, msg)
			})
		}
	}()
	return eg.Wait()
}

func (u *outbounds) Deliver(ctx context.Context, dest config.NodeID, msg interfaces.Message) error {
	var unicast *outbound
	func() {
		u.m.Lock()
		defer u.m.Unlock()

		unicast = u.table[dest]
	}()
	if unicast == nil {
		return interfaces.ErrUnknownPeer
	}
	return unicast.Deliver(ctx, msg)
}

func (u *outbound) Deliver(ctx context.Context, msg interfaces.Message) error {
	var list = make([]interfaces.RouterShink, 0, 16)
	func() {
		u.m.Lock()
		defer u.m.Unlock()

		for _, f := range u.table {
			list = append(list, f)
		}
	}()
	if len(list) == 0 {
		return interfaces.ErrPeerNoConnection
	}

	// Damn simple load balancing
	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

	var lastError error
	for _, f := range list {
		err := f(ctx, msg)
		if err == nil {
			return nil // Succes!
		}

		if lastError != nil {
			logger.GetFrom(ctx).Warnf("failed to deliver message, trying another connection...", "err", err)
		}
		lastError = err
	}
	return lastError
}
