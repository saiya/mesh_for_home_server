package peering

import (
	"context"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
	"github.com/saiya/mesh_for_home_server/tlshelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type peeringClient struct {
	router      interfaces.Router
	addr        string
	connections int

	ctx       context.Context
	ctxCancel context.CancelFunc

	gc *grpc.ClientConn

	stat interfaces.PeeringClientStat
}

func NewPeeringClient(parentCtx context.Context, cfg *config.PeeringConnectConfig, router interfaces.Router) (interfaces.PeeringClient, error) {
	c := peeringClient{
		router: router,
		addr:   cfg.Address,
	}
	c.ctx, c.ctxCancel = context.WithCancel(parentCtx)
	c.ctx = logger.Wrap(c.ctx, "addr", c.addr)
	logger.GetFrom(c.ctx).Debugw("Initializing peering clients...")

	err := newGRPCClient(cfg, &c)
	if err != nil {
		return nil, err
	}

	c.connections = cfg.Connections
	if c.connections <= 0 {
		c.connections = config.PeeringConnectionDefaultCount
	}
	retryInverval := time.Second * time.Duration(math.Max(3, float64(cfg.ConnectionRetryIntervalSec)))

	for connIdx := 0; connIdx < c.connections; connIdx++ {
		connCtx := logger.Wrap(c.ctx, "mutiplex", fmt.Sprintf("%d/%d", connIdx+1, c.connections))
		go func() {
		peering:
			for {
				conn := &peeringClientConnection{
					client: &c,
					router: c.router,
					ctx:    connCtx,
					stat:   &c.stat,
				}
				err := conn.peering()
				logger.GetFrom(conn.ctx).Infow(
					"Peering connection down, retry after "+retryInverval.String(),
					"err", err,
					"addr", c.addr,
				)

				retryTimer := time.NewTimer(retryInverval)
				select {
				case <-c.ctx.Done():
					break peering
				case <-retryTimer.C:
					continue
				}
			}
			logger.GetFrom(connCtx).Debugw("Peering connection ended")
		}()
	}
	return &c, nil
}

type peeringClientConnection struct {
	client *peeringClient
	router interfaces.Router
	ctx    context.Context

	stat *interfaces.PeeringClientStat
}

func (c *peeringClientConnection) peering() error {
	ctx := c.ctx
	logger.GetFrom(ctx).Debugw("Establishing peering connection...")

	atomic.AddUint64(&c.stat.PeeringAttempts, 1)
	var connSendLock sync.Mutex
	conn, err := generated.NewPeeringClient(c.client.gc).Peer(ctx)
	if err != nil {
		return err
	}
	defer func() {
		connSendLock.Lock()
		defer connSendLock.Unlock()
		debugLogIfErr(ctx, conn.CloseSend())
	}()
	atomic.AddUint64(&c.stat.PeeringConnected, 1)
	ctx = withConnectionLogAttributes(conn.Context())

	atomic.AddUint64(&c.stat.HandshakeAttempts, 1)
	handShakeResult, err := c.doHandshake(ctx, conn)
	if err != nil {
		logger.GetFrom(ctx)
		return err
	}
	atomic.AddUint64(&c.stat.HandshakeSucceeded, 1)
	logger.GetFrom(ctx).Debugw("Handshake done", "node", handShakeResult.peerNodeID)

	// FIXME: This code sending same message over ALL peering connections (= sending duplicated message)
	// Should implement living-connection list and round-robin
	sinkLoggingCtx := ctx
	deregisterShink := c.router.RegisterSink(handShakeResult.peerNodeID, func(ctx context.Context, msg interfaces.Message) error {
		logger.GetFrom(sinkLoggingCtx).Debugw("Sending peer message", "peer-msg", interfaces.MsgLogString(msg))

		connSendLock.Lock()
		defer connSendLock.Unlock()
		return conn.Send(&generated.PeerClientMessage{
			Message: &generated.PeerClientMessage_PeerMessage{
				PeerMessage: msg,
			},
		})
	})
	defer deregisterShink()

	for {
		packet, err := conn.Recv()
		if err != nil {
			logger.GetFrom(ctx).Debugw("Recv() failed, stream closed?", "err", err)
			return err
		}

		switch msg := packet.Message.(type) {
		case *generated.PeerServerMessage_CloseByServer:
			return fmt.Errorf("Stream closed by server, likely protocol error happened: %v", msg.CloseByServer.Reason)
		case *generated.PeerServerMessage_PeerMessage:
			atomic.AddUint64(&c.stat.PeerMessageReceived, 1)
			logger.GetFrom(ctx).Debugw("Peer message received", "peer-msg", interfaces.MsgLogString(msg.PeerMessage))
			c.router.Deliver(ctx, handShakeResult.peerNodeID, c.router.NodeID(), msg.PeerMessage)
		default:
			return fmt.Errorf("Unknown message received, closing connection")
		}
	}
}

func (c *peeringClient) String() string {
	return fmt.Sprintf("PeeringClient{ addr: %s, connections: %d, stat: %s }", c.addr, c.connections, c.Stat().String())
}

func (c *peeringClient) Close(ctx context.Context) error {
	c.ctxCancel()
	return c.gc.Close()
}

func (c *peeringClient) Stat() interfaces.PeeringClientStat {
	return c.stat.Clone()
}

func newGRPCClient(config *config.PeeringConnectConfig, c *peeringClient) error {
	var err error
	options := make([]grpc.DialOption, 0, 2)

	options = append(options, grpc.WithKeepaliveParams(
		keepalive.ClientParameters{
			Time:                time.Second * 15,
			PermitWithoutStream: true,
		},
	))

	var creds credentials.TransportCredentials
	if config.TLS != nil {
		tls, err := tlshelper.BuildTLSClientConfig(config.TLS)
		if err != nil {
			return err
		}
		creds = credentials.NewTLS(tls)
	} else {
		creds = insecure.NewCredentials()
	}
	options = append(options, grpc.WithTransportCredentials(creds))

	c.gc, err = grpc.Dial(config.Address, options...)
	if err != nil {
		return err
	}
	return nil
}
