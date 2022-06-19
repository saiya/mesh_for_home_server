package handler

import (
	"context"
	"io"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type pingHandler struct {
	router interfaces.Router
	listen io.Closer
}

func newPingHandler(router interfaces.Router) *pingHandler {
	h := pingHandler{
		router: router,
	}

	h.listen = router.Listen(func(ctx context.Context, from config.NodeID, msg interface{}) error {
		ping, ok := msg.(*generated.PeerMessage_Ping)
		if !ok {
			return nil
		}

		logger.GetFrom(ctx).Debugw("received PING, responding PONG", "payload", ping.Ping.Payload)
		return h.router.Deliver(ctx, from, &generated.PeerMessage_Pong{
			Pong: &generated.Pong{Payload: ping.Ping.Payload},
		})
	})
	return &h
}

func (h *pingHandler) Close(ctx context.Context) error {
	h.listen.Close()
	return nil
}
