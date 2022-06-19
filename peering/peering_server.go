package peering

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
	"github.com/saiya/mesh_for_home_server/tlshelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type peeringServer struct {
	nodeID config.NodeID

	listenAddr  string
	tcpListener net.Listener

	port int
	gs   *grpc.Server

	closeWait *sync.WaitGroup
	closeErr  error

	stat interfaces.PeeringServerStat
}

func NewPeeringServer(nodeID config.NodeID, config *config.PeeringAcceptConfig) (interfaces.PeeringServer, error) {
	s := peeringServer{
		nodeID:    nodeID,
		closeWait: &sync.WaitGroup{},
	}
	err := newListener(config, &s)
	if err != nil {
		return nil, err
	}

	s.gs, err = newGRPCServer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to start gRPC server: %w", err)
	}

	generated.RegisterPeeringServer(s.gs, &s)

	s.closeWait.Add(1)
	go func() {
		defer s.closeWait.Done()

		logger.Get().Infow("Peering server listening at "+s.listenAddr, "port", s.port)
		s.closeErr = s.gs.Serve(s.tcpListener)
	}()
	return &s, nil
}

func newListener(config *config.PeeringAcceptConfig, s *peeringServer) error {
	var err error
	if config.Listen != "" {
		s.listenAddr = config.Listen
	} else {
		s.listenAddr = ":0"
	}

	s.tcpListener, err = net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen \"%s\": %w", s.listenAddr, err)
	}

	s.port = s.tcpListener.Addr().(*net.TCPAddr).Port
	return nil
}

func newGRPCServer(config *config.PeeringAcceptConfig) (*grpc.Server, error) {
	options := make([]grpc.ServerOption, 0, 2)
	options = append(options, grpc.KeepaliveEnforcementPolicy(
		keepalive.EnforcementPolicy{
			PermitWithoutStream: true,
		},
	))
	if config.TLS != nil {
		tlsConfig, err := tlshelper.BuildTLSServerConfig(config.TLS)
		if err != nil {
			return nil, fmt.Errorf("failed to load gRPC server TLS cert and key: %w", err)
		}
		options = append(options, grpc.Creds(credentials.NewTLS(tlsConfig)))
	}
	return grpc.NewServer(options...), nil
}

func (s *peeringServer) Close(ctx context.Context) error {
	s.gs.Stop()
	s.closeWait.Wait()
	return s.closeErr
}

func (s *peeringServer) String() string {
	return fmt.Sprintf(
		"PeeringServer{ grpc_port: %d, stat: %s }",
		s.port,
		s.Stat().String(),
	)
}

func (s *peeringServer) Stat() interfaces.PeeringServerStat {
	return s.stat
}

func (s *peeringServer) Peer(conn generated.Peering_PeerServer) error {
	atomic.AddUint64(&s.stat.PeeringConnections, 1)
	sm := newStateMachine(
		conn.Context(), true,
		func(msg *generated.PeerMessage) error {
			return conn.Send(&generated.PeerServerMessage{Message: &generated.PeerServerMessage_PeerMessage{PeerMessage: msg}})
		},
	)

	atomic.AddUint64(&s.stat.HandshakeAttempts, 1)
	if err := s.doHandshake(conn, sm); err != nil {
		sm.Abort(generated.PeeringAbort_HANDSHAKE_FAILURE, err)
		return nil
	}
	atomic.AddUint64(&s.stat.HandshakeSucceeded, 1)

	go func() {
		for sm.Alive() {
			packet, err := conn.Recv()
			if err != nil {
				sm.Abort(generated.PeeringAbort_STREAM_CLOSED, err)
				break
			}

			switch msg := packet.Message.(type) {
			case *generated.PeerClientMessage_PeerMessage:
				atomic.AddUint64(&s.stat.PeerMessageReceived, 1)
				sm.Update(msg.PeerMessage)
			default:
				sm.Abort(generated.PeeringAbort_INVALID_MESSAGE_ORDER, nil)
			}
		}
	}()
	_, _ = <-sm.ctx.Done()
	return nil
}
