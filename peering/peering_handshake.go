package peering

import (
	"context"
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type handshakeResult struct {
	peerNodeID config.NodeID
}

func (s *peeringServer) doHandshake(ctx context.Context, conn generated.Peering_PeerServer) (handshakeResult, error) {
	logger.GetFrom(ctx).Debugw("Start of handshake...")
	result := handshakeResult{}

	helloMsg, err := conn.Recv()
	if err != nil {
		return result, fmt.Errorf("failed to receive CLIENT HELLO message: %w", err)
	}
	hello, ok := helloMsg.Message.(*generated.PeerClientMessage_ClientHello)
	if !ok {
		return result, fmt.Errorf("expected CLIENT HELLO message but got different message")
	}
	result.peerNodeID = config.NodeID(hello.ClientHello.GetNodeId())

	err = conn.Send(&generated.PeerServerMessage{
		Message: &generated.PeerServerMessage_ServerHello{
			ServerHello: &generated.ServerHello{
				NodeId: string(s.nodeID),
			},
		},
	})
	if err != nil {
		return result, fmt.Errorf("failed to send SERVER HELLO: %w", err)
	}

	answerMsg, err := conn.Recv()
	if err != nil {
		return result, fmt.Errorf("failed to receive HANDSHAKE DONE message: %w", err)
	}
	_, ok = answerMsg.Message.(*generated.PeerClientMessage_HandshakeDone)
	if !ok {
		return result, fmt.Errorf("expected HANDSHAKE DONE message but got different message")
	}

	logger.GetFrom(ctx).Debugw("Handshake done")
	return result, nil
}

func (c *peeringClient) doHandshake(ctx context.Context, conn generated.Peering_PeerClient) (handshakeResult, error) {
	logger.GetFrom(ctx).Debugw("Start of handshake...")
	result := handshakeResult{}

	err := conn.Send(&generated.PeerClientMessage{
		Message: &generated.PeerClientMessage_ClientHello{
			ClientHello: &generated.ClientHello{
				NodeId: string(c.router.NodeID()),
			},
		},
	})
	if err != nil {
		return result, fmt.Errorf("failed to send CLIENT HELLO: %w", err)
	}

	helloMsg, err := conn.Recv()
	if err != nil {
		return result, fmt.Errorf("failed to receive SERVER HELLO message: %w", err)
	}
	hello, ok := helloMsg.Message.(*generated.PeerServerMessage_ServerHello)
	if !ok {
		return result, fmt.Errorf("expected SERVER HELLO message but got different message")
	}
	result.peerNodeID = config.NodeID(hello.ServerHello.NodeId)

	err = conn.Send(&generated.PeerClientMessage{
		Message: &generated.PeerClientMessage_HandshakeDone{
			HandshakeDone: &generated.HandShakeDone{},
		},
	})
	if err != nil {
		return result, fmt.Errorf("failed to send HANDSHAKE DONE: %w", err)
	}

	logger.GetFrom(ctx).Debugw("Handshake done")
	return result, nil
}
