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
	nodeID config.NodeID
	addr   string

	ctx       context.Context
	ctxCancel context.CancelFunc

	gc *grpc.ClientConn

	stat interfaces.PeeringClientStat
}

func NewPeeringClient(parentCtx context.Context, nodeID config.NodeID, config *config.PeeringConnectConfig) (interfaces.PeeringClient, error) {
	c := peeringClient{
		nodeID: nodeID,
		addr:   config.Address,
	}
	c.ctx, c.ctxCancel = context.WithCancel(parentCtx)

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
	atomic.AddUint64(&c.stat.PeeringAttempts, 1)
	conn, err := generated.NewPeeringClient(c.gc).Peer(c.ctx)
	if err != nil {
		return err
	}
	atomic.AddUint64(&c.stat.PeeringConnected, 1)

	sm := newStateMachine(
		c.ctx, false,
		func(msg *generated.PeerMessage) error {
			return conn.Send(&generated.PeerClientMessage{Message: &generated.PeerClientMessage_PeerMessage{PeerMessage: msg}})
		},
	)

	atomic.AddUint64(&c.stat.HandshakeAttempts, 1)
	err = c.doHandshake(conn, sm)
	if err != nil {
		return err
	}
	atomic.AddUint64(&c.stat.HandshakeSucceeded, 1)

	for sm.Alive() {
		packet, err := conn.Recv()
		if err != nil {
			sm.Abort(generated.PeeringAbort_STREAM_CLOSED, err)
			return err
		}

		switch msg := packet.Message.(type) {
		case *generated.PeerServerMessage_PeerMessage:
			atomic.AddUint64(&c.stat.PeerMessageReceived, 1)
			sm.Update(msg.PeerMessage)
		default:
			sm.Abort(generated.PeeringAbort_INVALID_MESSAGE_ORDER, nil)
			return fmt.Errorf("PeeringAbort_INVALID_MESSAGE_ORDER")
		}
	}
	return fmt.Errorf("Stream dead")
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
