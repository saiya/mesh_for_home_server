package egress

import (
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/egress/handler"
	"github.com/saiya/mesh_for_home_server/egress/handler/httphandler"
	"github.com/saiya/mesh_for_home_server/interfaces"
)

func StartEgress(c *config.EgressConfigs, router interfaces.Router) ([]interfaces.MessageHandler, error) {
	httpHandler := httphandler.NewHttpHandler(router)
	handlers := []interfaces.MessageHandler{
		handler.NewPingHandler(router),
		httpHandler,
	}

	if c == nil {
		return handlers, nil
	}

	for i := range c.HTTP {
		if err := httpHandler.AddEgress(&c.HTTP[i]); err != nil {
			return handlers, fmt.Errorf("failed to start HTTP(S) egress [%d]: %w", i, err)
		}
	}

	return handlers, nil
}
