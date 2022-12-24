package peering

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/egress/handler"
	"github.com/saiya/mesh_for_home_server/ingress/forwarder"
	"github.com/saiya/mesh_for_home_server/router"
	"github.com/saiya/mesh_for_home_server/tlshelper/tlstesting"
	"github.com/stretchr/testify/assert"
)

func TestPeering(t *testing.T) {
	ctx := context.Background()

	serverRouter := router.NewRouter("server-")
	serverPingHandler := handler.NewPingHandler(serverRouter)
	defer serverPingHandler.Close(ctx)
	pingFromServer := forwarder.NewPingForwarder(serverRouter)
	defer pingFromServer.Close(ctx)

	clientRouter := router.NewRouter("client-")
	clientPingHandler := handler.NewPingHandler(clientRouter)
	defer clientPingHandler.Close(ctx)
	pingFromClient := forwarder.NewPingForwarder(clientRouter)
	defer pingFromClient.Close(ctx)

	serverTLSConfig := tlstesting.GenerateServerCert("localhost")
	clientTLSConfig := tlstesting.EnableClientCert(serverTLSConfig)

	server, err := NewPeeringServer(
		&config.PeeringAcceptConfig{
			Listen: "localhost:0",
			TLS:    serverTLSConfig,
		},
		serverRouter,
	)
	assert.NoError(t, err)
	defer server.Close(ctx)

	assert.Equal(t, uint64(0), server.Stat().HandshakeSucceeded)

	client, err := NewPeeringClient(
		ctx,
		&config.PeeringConnectConfig{
			Address: fmt.Sprintf("localhost:%d", server.(*peeringServer).port),
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
