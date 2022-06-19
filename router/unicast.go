package router

import (
	"context"
	"math/rand"
	"sync"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
)

type unicasts struct {
	m     sync.Mutex
	table map[config.NodeID]*unicast
}

type unicast struct {
	m           sync.Mutex
	idGenerator uint64
	table       map[uint64]interfaces.RouterShink
}

func newUnicasts() *unicasts {
	return &unicasts{
		table: make(map[config.NodeID]*unicast),
	}
}

func newUnicast() *unicast {
	return &unicast{
		table: make(map[uint64]interfaces.RouterShink),
	}
}

func (u *unicasts) Register(dest config.NodeID, callback interfaces.RouterShink) interfaces.RouterUnregister {
	var unicast *unicast
	func() {
		u.m.Lock()
		defer u.m.Unlock()

		unicast = u.table[dest]
		if unicast == nil {
			unicast = newUnicast()
			u.table[dest] = unicast
		}
	}()
	return unicast.Register(callback)
}

func (u *unicast) Register(callback interfaces.RouterShink) interfaces.RouterUnregister {
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

func (u *unicasts) Deliver(ctx context.Context, dest config.NodeID, msg interfaces.Message) error {
	var unicast *unicast
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

func (u *unicast) Deliver(ctx context.Context, msg interfaces.Message) error {
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
