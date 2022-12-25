package integration_test_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/egress/handler"
	"github.com/saiya/mesh_for_home_server/ingress/forwarder"
	"github.com/saiya/mesh_for_home_server/peering"
	"github.com/saiya/mesh_for_home_server/router"
	"github.com/saiya/mesh_for_home_server/tlshelper/tlstesting"
	"github.com/stretchr/testify/assert"
)

func TestPeering(t *testing.T) {
	ctx := context.Background()

	serverRouter := router.NewRouter("server")
	defer serverRouter.Close(ctx)
	serverPingHandler := handler.NewPingHandler(serverRouter)
	defer serverPingHandler.Close(ctx)
	pingFromServer := forwarder.NewPingForwarder(serverRouter)
	defer pingFromServer.Close(ctx)

	clientRouter := router.NewRouter("client")
	defer clientRouter.Close(ctx)
	clientPingHandler := handler.NewPingHandler(clientRouter)
	defer clientPingHandler.Close(ctx)
	pingFromClient := forwarder.NewPingForwarder(clientRouter)
	defer pingFromClient.Close(ctx)

	serverTLSConfig := tlstesting.GenerateServerCert("localhost")
	clientTLSConfig := tlstesting.EnableClientCert(serverTLSConfig)

	server, err := peering.NewPeeringServer(
		&config.PeeringAcceptConfig{
			Listen: "localhost:0",
			TLS:    serverTLSConfig,
		},
		serverRouter,
	)
	assert.NoError(t, err)
	defer server.Close(ctx)

	assert.Equal(t, uint64(0), server.Stat().HandshakeSucceeded)

	client, err := peering.NewPeeringClient(
		ctx,
		&config.PeeringConnectConfig{
			Address: fmt.Sprintf("localhost:%d", server.Port()),
			TLS:     clientTLSConfig,
		},
		clientRouter,
	)
	assert.NoError(t, err)
	defer client.Close(ctx)

	isConnected := func() bool {
		return (server.Stat().HandshakeSucceeded > 0) && (client.Stat().HandshakeSucceeded > 0)
	}
	waitUntil := time.Now().Add(time.Second * 3)
	for time.Now().Before(waitUntil) {
		if isConnected() {
			break
		}
	}
	assert.True(t, isConnected())

	assert.NoError(t, pingFromServer.Ping(ctx, clientRouter.NodeID(), time.Second*3))
	assert.NoError(t, pingFromClient.Ping(ctx, serverRouter.NodeID(), time.Second*3))
}
