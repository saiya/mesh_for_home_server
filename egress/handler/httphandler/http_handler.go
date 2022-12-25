package httphandler

import (
	"context"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type HttpHandler interface {
	interfaces.MessageHandler
	AddEgress(c *config.HTTPEgressConfig) error
	Advertise() *generated.HttpAdvertisement
}
type httpHandler struct{}

func NewHttpHandler(router interfaces.Router) HttpHandler {
	// FIXME: Implement
	return &httpHandler{}
}

func (h *httpHandler) Advertise() *generated.HttpAdvertisement {
	// FIXME: Implement
	return nil
}

func (h *httpHandler) AddEgress(c *config.HTTPEgressConfig) error {
	// FIXME: Implement
	return nil
}

func (h *httpHandler) Close(ctx context.Context) error {
	// FIXME: Implement
	return nil
}
