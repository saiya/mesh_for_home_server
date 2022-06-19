package peering

import (
	"context"
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
	"github.com/saiya/mesh_for_home_server/tlshelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type peeringClient struct {
	router interfaces.Router
	addr   string

	ctx       context.Context
	ctxCancel context.CancelFunc

	gc *grpc.ClientConn

	stat interfaces.PeeringClientStat
}

func NewPeeringClient(parentCtx context.Context, config *config.PeeringConnectConfig, router interfaces.Router) (interfaces.PeeringClient, error) {
	c := peeringClient{
		router: router,
		addr:   config.Address,
	}
	c.ctx, c.ctxCancel = context.WithCancel(parentCtx)
	c.ctx = logger.Wrap(c.ctx, "addr", c.addr)

	err := newGRPCClient(config, &c)
	if err != nil {
		return nil, err
	}

	retryInverval := time.Second * time.Duration(math.Max(3, float64(config.ConnectionRetryIntervalSec)))
	go func() {
	peering:
		for {
			err := c.peering()
			logger.Get().Infow(
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
	}()

	return &c, nil
}

func (c *peeringClient) peering() error {
	ctx := c.ctx

	atomic.AddUint64(&c.stat.PeeringAttempts, 1)
	conn, err := generated.NewPeeringClient(c.gc).Peer(ctx)
	if err != nil {
		return err
	}
	defer conn.CloseSend()
	atomic.AddUint64(&c.stat.PeeringConnected, 1)
	ctx = withConnectionLogAttributes(conn.Context())

	atomic.AddUint64(&c.stat.HandshakeAttempts, 1)
	handShakeResult, err := c.doHandshake(ctx, conn)
	if err != nil {
		logger.GetFrom(ctx)
		return err
	}
	atomic.AddUint64(&c.stat.HandshakeSucceeded, 1)

	sinkRegistration := c.router.RegisterSink(handShakeResult.peerNodeID, func(ctx context.Context, msg interfaces.Message) error {
		logger.GetFrom(ctx).Debugw("Sending peer message", "peer-msg", interfaces.MsgLogString(msg))
		return conn.Send(&generated.PeerClientMessage{
			Message: &generated.PeerClientMessage_PeerMessage{
				PeerMessage: msg,
			},
		})
	})
	defer sinkRegistration.Close()

	for {
		packet, err := conn.Recv()
		if err != nil {
			logger.GetFrom(ctx).Debugw("Recv() failed, stream closed?", "err", err)
			return err
		}

		switch msg := packet.Message.(type) {
		case *generated.PeerServerMessage_CloseByServer:
			return fmt.Errorf("Stream closed by server, likely protocol error happend: %v", msg.CloseByServer.Reason)
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
	return fmt.Sprintf("PeeringClient{ addr: %s, stat: %s }", c.addr, c.Stat().String())
}

func (c *peeringClient) Close(ctx context.Context) error {
	c.ctxCancel()
	return c.gc.Close()
}

func (c *peeringClient) Stat() interfaces.PeeringClientStat {
	return c.stat
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

	if config.TLS != nil {
		tls, err := tlshelper.BuildTLSClientConfig(config.TLS)
		if err != nil {
			return err
		}
		options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(tls)))
	}

	c.gc, err = grpc.Dial(config.Address, options...)
	if err != nil {
		return err
	}
	return nil
}
