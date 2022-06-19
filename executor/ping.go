package executor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type pingExecutor struct {
	router           interfaces.Router
	unregisterListen interfaces.RouterUnregister

	m     sync.Mutex
	table map[string]*inFlightPing
}

type inFlightPing struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	sucessCh  chan struct{}
}

func NewPingExecutor(router interfaces.Router) interfaces.PingExecutor {
	exec := &pingExecutor{
		router: router,
		table:  make(map[string]*inFlightPing),
	}
	exec.unregisterListen = router.Listen(func(ctx context.Context, from config.NodeID, msg interfaces.Message) error {
		pong := msg.GetPong()
		if pong != nil {
			exec.Pong(ctx, pong)
		}
		return nil
	})
	return exec
}

func (exec *pingExecutor) Close(ctx context.Context) error {
	exec.m.Lock()
	defer exec.m.Unlock()

	exec.unregisterListen()
	return nil
}

func (exec *pingExecutor) Ping(ctx context.Context, dest config.NodeID, timeout time.Duration) error {
	nonce := exec.generateNonce(dest)
	logger.GetFrom(ctx).Infow("Executing PING...", "dest", dest, "timeout", timeout, "nonce", nonce)

	inflight := &inFlightPing{
		sucessCh: make(chan struct{}),
	}
	inflight.ctx, inflight.ctxCancel = context.WithTimeout(ctx, timeout)
	defer inflight.ctxCancel()

	func() {
		exec.m.Lock()
		defer exec.m.Unlock()
		exec.table[nonce] = inflight
	}()
	defer func() {
		exec.m.Lock()
		defer exec.m.Unlock()
		delete(exec.table, nonce)
	}()

	exec.router.Deliver(ctx, exec.router.NodeID(), dest, &generated.PeerMessage{
		Message: &generated.PeerMessage_Ping{
			Ping: &generated.Ping{
				Payload: nonce,
			},
		},
	})

	select {
	case <-inflight.sucessCh:
		return nil
	case <-inflight.ctx.Done():
		return inflight.ctx.Err()
	}
}

func (exec *pingExecutor) Pong(ctx context.Context, msg *generated.Pong) {
	logger.GetFrom(ctx).Infow("PONG received", "nonce", msg.Payload)

	var inflight *inFlightPing
	func() {
		exec.m.Lock()
		defer exec.m.Unlock()
		inflight = exec.table[msg.Payload]
	}()
	if inflight == nil {
		logger.GetFrom(ctx).Warnw("PONG received but unknown value, skipping...", "nonce", msg.Payload)
	}

	close(inflight.sucessCh)
}

func (exec *pingExecutor) generateNonce(dest config.NodeID) string {
	return fmt.Sprintf("%s$%s$nonce:%s", exec.router.NodeID(), dest, uuid.New().String())
}
