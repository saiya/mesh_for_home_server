package peering

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/tlshelper/tlstesting"
	"github.com/stretchr/testify/assert"
)

func TestPeering(t *testing.T) {
	logger.EnableDebugLog()

	ctx := context.Background()
	serverTLSConfig := tlstesting.GenerateServerCert("localhost")
	clientTLSConfig := tlstesting.EnableClientCert(serverTLSConfig)

	server, err := NewPeeringServer(
		config.NodeID("test-server"),
		&config.PeeringAcceptConfig{
			Listen: "localhost:0",
			TLS:    serverTLSConfig,
		},
	)
	assert.NoError(t, err)
	defer server.Close(ctx)

	assert.Equal(t, uint64(0), server.Stat().HandshakeSucceeded)

	client, err := NewPeeringClient(
		ctx, config.NodeID("test-client"),
		&config.PeeringConnectConfig{
			Address: fmt.Sprintf("localhost:%d", server.(*peeringServer).port),
			TLS:     clientTLSConfig,
		},
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
}
