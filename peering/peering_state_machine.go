package peering

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type peeringStateMachine struct {
	server       bool
	connectionID uint64 // This ID is local in this node, not globally unique
	emit         func(*generated.PeerMessage) error

	mu        sync.Mutex
	ctx       context.Context
	ctxCancel context.CancelFunc

	aborted    bool
	peerNodeID config.NodeID
}

var connectionIDGen = uint64(0)

func newStateMachine(
	parentContext context.Context,
	server bool,
	emit func(*generated.PeerMessage) error,
) *peeringStateMachine {
	sm := &peeringStateMachine{
		server:       server,
		connectionID: atomic.AddUint64(&connectionIDGen, 1),
		emit:         emit,
	}
	sm.ctx, sm.ctxCancel = context.WithCancel(parentContext)
	return sm
}

func (sm *peeringStateMachine) PeerNodeID() config.NodeID {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return sm.peerNodeID
}

func (sm *peeringStateMachine) SetPeerNodeID(peerNodeID config.NodeID) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.peerNodeID = peerNodeID
}

func (sm *peeringStateMachine) Alive() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return !sm.aborted
}

func (sm *peeringStateMachine) Abort(code generated.PeeringAbort_PeeringError, err error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.aborted {
		return
	}

	defer sm.goToAbortedState()

	emitResult := sm.emit(&generated.PeerMessage{
		Message: &generated.PeerMessage_Abort{
			Abort: &generated.PeeringAbort{
				Error: code,
			},
		},
	})
	logger.Get().Infow(
		"Emitted abort message",
		"err", err, "abort-emit-result", emitResult,
		"connection-id", sm.connectionID,
		"peer-node-id", sm.peerNodeID,
		"code", code,
	)
}

func (sm *peeringStateMachine) goToAbortedState() {
	sm.aborted = true
	sm.ctxCancel()
}

func (sm *peeringStateMachine) Update(msgWrapper *generated.PeerMessage) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	logger.Get().Debugw(
		"Received peer-message",
		"connection-id", sm.connectionID,
		"peer-node-id", sm.peerNodeID,
		"msg", msgWrapper,
	)

	if sm.aborted {
		return
	}

	switch msg := msgWrapper.Message.(type) {
	case *generated.PeerMessage_Abort:
		sm.goToAbortedState()
	default:
		// TODO: ... Implement ...
		logger.Get().Debugw("... Not yet implemented ...", "msg", msg)
	}
}
