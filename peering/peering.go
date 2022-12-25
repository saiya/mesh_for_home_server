package peering

import (
	"context"
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
)

func StartPeering(ctx context.Context, config *config.PeeringConfig, router interfaces.Router) (interfaces.PeeringServer, []interfaces.PeeringClient, error) {
	var server interfaces.PeeringServer
	clients := make([]interfaces.PeeringClient, 0)

	if config.Accept != nil {
		server, err := NewPeeringServer(config.Accept, router)
		if err != nil {
			return server, clients, fmt.Errorf("failed to start peering server: %w", err)
		}
	}

	for _, connect := range config.Connect {
		client, err := NewPeeringClient(ctx, connect, router)
		if err != nil {
			return server, clients, fmt.Errorf("failed to start peering server: %w", err)
		}
		clients = append(clients, client)
	}

	return server, clients, nil
}
