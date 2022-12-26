package peering

import (
	"context"
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
)

func StartPeering(ctx context.Context, config *config.PeeringConfig, router interfaces.Router) ([]interfaces.PeeringServer, []interfaces.PeeringClient, error) {
	servers := make([]interfaces.PeeringServer, 0, 1) // To avoid nil interface handling issue, we return it as array
	clients := make([]interfaces.PeeringClient, 0)

	if config.Accept != nil {
		server, err := NewPeeringServer(config.Accept, router)
		if err != nil {
			return servers, clients, fmt.Errorf("failed to start peering server: %w", err)
		}
		servers = append(servers, server)
	}

	for _, connect := range config.Connect {
		client, err := NewPeeringClient(ctx, connect, router)
		if err != nil {
			return servers, clients, fmt.Errorf("failed to start peering server: %w", err)
		}
		clients = append(clients, client)
	}

	return servers, clients, nil
}
