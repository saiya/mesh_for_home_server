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
	router interfaces.Router

	listenAddr  string
	tcpListener net.Listener

	port int
	gs   *grpc.Server

	closeWait *sync.WaitGroup
	closeErr  error

	stat interfaces.PeeringServerStat
}

func NewPeeringServer(config *config.PeeringAcceptConfig, router interfaces.Router) (interfaces.PeeringServer, error) {
	s := peeringServer{
		router:    router,
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
	return s.stat.Clone()
}

func (s *peeringServer) Peer(conn generated.Peering_PeerServer) error {
	atomic.AddUint64(&s.stat.PeeringConnections, 1)
	ctx := withConnectionLogAttributes(conn.Context())

	atomic.AddUint64(&s.stat.HandshakeAttempts, 1)
	handshakeResult, err := s.doHandshake(ctx, conn)
	if err != nil {
		closeByServer(ctx, conn, generated.CloseByServerReason_HANDSHAKE_FAILURE)
		return nil
	}
	if handshakeResult.peerNodeID == s.router.NodeID() {
		logger.GetFrom(ctx).Warnw("NodeID duplication: peer's nodeID equals with this server's nodeID", "node-id", s.router.NodeID())
		closeByServer(ctx, conn, generated.CloseByServerReason_HANDSHAKE_FAILURE)
		return nil
	}
	atomic.AddUint64(&s.stat.HandshakeSucceeded, 1)

	deregisterShink := s.router.RegisterSink(handshakeResult.peerNodeID, func(parentCtx context.Context, msg interfaces.Message) error {
		ctx := withConnectionLogAttributes(parentCtx)
		logger.GetFrom(ctx).Debugw("Sending peer message", "peer-msg", interfaces.MsgLogString(msg))
		return conn.Send(
			&generated.PeerServerMessage{
				Message: &generated.PeerServerMessage_PeerMessage{
					PeerMessage: msg,
				},
			},
		)
	})
	defer deregisterShink()

	go func() {
		for {
			packet, err := conn.Recv()
			if err != nil {
				logger.GetFrom(ctx).Infow("Recv() failed, connection closed?", "err", err)
				closeByServer(ctx, conn, generated.CloseByServerReason_CONNECTION_DOWN)
				break
			}

			switch msg := packet.Message.(type) {
			case *generated.PeerClientMessage_Routing:
				atomic.AddUint64(&s.stat.RouteRecived, 1)
				s.router.Update(ctx, msg.Routing)
			case *generated.PeerClientMessage_PeerMessage:
				atomic.AddUint64(&s.stat.PeerMessageReceived, 1)
				s.router.Deliver(ctx, handshakeResult.peerNodeID, s.router.NodeID(), msg.PeerMessage)
			default:
				logger.GetFrom(ctx).Warnw("Unknown message received, closing connection")
				closeByServer(ctx, conn, generated.CloseByServerReason_INVALID_MESSAGE)
			}
		}
	}()
	<-ctx.Done()
	return nil
}

func closeByServer(ctx context.Context, conn generated.Peering_PeerServer, reason generated.CloseByServerReason) {
	err := conn.Send(
		&generated.PeerServerMessage{
			Message: &generated.PeerServerMessage_CloseByServer{
				CloseByServer: &generated.PeerCloseByServer{
					Reason: reason,
				},
			},
		},
	)
	logger.GetFrom(ctx).Infow("Closed connection by this server", "reason", reason, "close-message-send-result", err)
}
