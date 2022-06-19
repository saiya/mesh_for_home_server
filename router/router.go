package router

import (
	"context"
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"

	"github.com/google/uuid"
)

type router struct {
	nodeID config.NodeID

	unicast   *unicasts
	listeners *broadcast
}

func NewRouter(nodeIDPrefix string) interfaces.Router {
	router := &router{
		nodeID: config.NodeID(fmt.Sprintf("%s%s", nodeIDPrefix, uuid.New().String())),

		unicast:   newUnicasts(),
		listeners: newBroardcast(),
	}
	return router
}

func (r *router) NodeID() config.NodeID {
	return r.nodeID
}

func (r *router) RegisterSink(dest config.NodeID, callback interfaces.RouterShink) interfaces.RouterUnregister {
	if r.nodeID == dest {
		panic("Sink's destination nodeID must be different from this node's ID")
	}

	return r.unicast.Register(dest, callback)
}

func (r *router) Listen(callback interfaces.RouterListener) interfaces.RouterUnregister {
	return r.listeners.Register(callback)
}

func (r *router) Deliver(parentCtx context.Context, from config.NodeID, dest config.NodeID, msg interfaces.Message) {
	ctx := logger.Wrap(parentCtx, "from", from, "dest", dest, "peer-msg", interfaces.MsgLogString(msg))
	if r.nodeID == dest {
		r.listeners.Invoke(ctx, from, msg)
	} else {
		err := r.unicast.Deliver(ctx, dest, msg)
		if err != nil {
			logger.GetFrom(ctx).Warnw("failed to deliver message", "err", err)
		}
	}
}
