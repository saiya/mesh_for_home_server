package peering

import (
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

func (s *peeringServer) doHandshake(conn generated.Peering_PeerServer, sm *peeringStateMachine) error {
	helloMsg, err := conn.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive CLIENT HELLO message: %w", err)
	}
	hello, ok := helloMsg.Message.(*generated.PeerClientMessage_ClientHello)
	if !ok {
		return fmt.Errorf("expected CLIENT HELLO message but got different message")
	}
	sm.SetPeerNodeID(config.NodeID(hello.ClientHello.GetNodeId()))

	err = conn.Send(&generated.PeerServerMessage{
		Message: &generated.PeerServerMessage_ServerHello{
			ServerHello: &generated.ServerHello{
				NodeId: string(s.nodeID),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send SERVER HELLO: %w", err)
	}

	answerMsg, err := conn.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive HANDSHAKE DONE message: %w", err)
	}
	_, ok = answerMsg.Message.(*generated.PeerClientMessage_HandshakeDone)
	if !ok {
		return fmt.Errorf("expected HANDSHAKE DONE message but got different message")
	}
	return nil
}

func (c *peeringClient) doHandshake(conn generated.Peering_PeerClient, sm *peeringStateMachine) error {
	err := conn.Send(&generated.PeerClientMessage{
		Message: &generated.PeerClientMessage_ClientHello{
			ClientHello: &generated.ClientHello{
				NodeId: string(c.nodeID),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send CLIENT HELLO: %w", err)
	}

	helloMsg, err := conn.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive SERVER HELLO message: %w", err)
	}
	hello, ok := helloMsg.Message.(*generated.PeerServerMessage_ServerHello)
	if !ok {
		return fmt.Errorf("expected SERVER HELLO message but got different message")
	}
	sm.SetPeerNodeID(config.NodeID(hello.ServerHello.NodeId))

	err = conn.Send(&generated.PeerClientMessage{
		Message: &generated.PeerClientMessage_HandshakeDone{
			HandshakeDone: &generated.HandShakeDone{},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send HANDSHAKE DONE: %w", err)
	}
	return nil
}
