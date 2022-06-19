package handler

import (
	"context"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type pingHandler struct {
	router    interfaces.Router
	endListen interfaces.RouterUnregister
}

func NewPingHandler(router interfaces.Router) interfaces.MessageHandler {
	h := pingHandler{
		router: router,
	}

	h.endListen = router.Listen(func(ctx context.Context, from config.NodeID, msg interfaces.Message) error {
		ping, ok := msg.Message.(*generated.PeerMessage_Ping)
		if !ok {
			return nil
		}

		logger.GetFrom(ctx).Debugw("received PING, responding PONG", "payload", ping.Ping.Payload)
		h.router.Deliver(ctx, router.NodeID(), from, &generated.PeerMessage{
			Message: &generated.PeerMessage_Pong{
				Pong: &generated.Pong{Payload: ping.Ping.Payload},
			},
		})
		return nil
	})
	return &h
}

func (h *pingHandler) Close(ctx context.Context) error {
	h.endListen()
	return nil
}
