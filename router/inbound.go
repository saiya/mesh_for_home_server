package router

import (
	"context"
	"sync"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
)

type inbound struct {
	m sync.Mutex

	listenerIDGen uint64
	table         map[uint64]interfaces.RouterListener
}

func newInbound() *inbound {
	return &inbound{
		table: make(map[uint64]interfaces.RouterListener),
	}
}

func (b *inbound) Invoke(parentCtx context.Context, from config.NodeID, msg interfaces.Message) {
	ctx := logger.Wrap(parentCtx, "from", from, "peer-msg", interfaces.MsgLogString(msg))

	for _, handler := range b.Snapshot() {
		err := handler(ctx, from, msg)
		if err != nil {
			logger.GetFrom(ctx).Errorf("messagehandler returned error: %w", err)
		}
	}
}

func (b *inbound) Snapshot() []interfaces.RouterListener {
	b.m.Lock()
	defer b.m.Unlock()

	snapshot := make([]interfaces.RouterListener, 0, len(b.table))
	for _, f := range b.table {
		snapshot = append(snapshot, f)
	}
	return snapshot
}

func (b *inbound) Register(f interfaces.RouterListener) interfaces.RouterUnregister {
	b.m.Lock()
	defer b.m.Unlock()

	b.listenerIDGen++
	id := b.listenerIDGen
	b.table[id] = f

	return func() {
		b.m.Lock()
		defer b.m.Unlock()

		delete(b.table, id)
	}
}
