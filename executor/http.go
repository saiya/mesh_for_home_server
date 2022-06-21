package executor

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type httpExecutor struct {
	router interfaces.Router

	m                sync.Mutex
	unregisterListen interfaces.RouterUnregister
	seqGen           uint64
	inflight         map[httpRequestID]*inFlightHTTP
}

type httpRequestID struct {
	destNodeID config.NodeID
	sequence   uint64
}

func NewHTTPExecutor(router interfaces.Router) interfaces.HTTPExecutor {
	exec := &httpExecutor{
		router:   router,
		inflight: make(map[httpRequestID]*inFlightHTTP),
	}
	exec.unregisterListen = router.Listen(exec.HandleResponse)
	return exec
}

func (exec *httpExecutor) Close(ctx context.Context) error {
	exec.m.Lock()
	defer exec.m.Unlock()

	exec.unregisterListen()
	for _, state := range exec.inflight {
		state.Close()
	}
	return nil
}

func (exec *httpExecutor) RoundTrip(req *http.Request) (*http.Response, error) {
	dest, err := exec.router.RouteHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	ctx := logger.Wrap(req.Context(), "dest", dest, "http_method", req.Method, "http_hostname", req.Host, "http_path", req.URL.Path)
	state := newInflightHTTP(req, dest)
	defer state.Close()

	id := httpRequestID{destNodeID: dest}
	func() {
		exec.m.Lock()
		defer exec.m.Unlock()

		exec.seqGen++
		id.sequence = exec.seqGen
		state.id = id

		exec.inflight[id] = state
	}()
	defer func() {
		exec.m.Lock()
		defer exec.m.Unlock()

		delete(exec.inflight, id)
	}()

	err = state.SendRequest(ctx, exec.router)
	if err != nil {
		return nil, err
	}
	return state.AwaitResponseHeader(ctx)
}

func (exec *httpExecutor) HandleResponse(ctx context.Context, from config.NodeID, msg interfaces.Message) error {
	switch msg.Message.(type) {
	case *generated.PeerMessage_HttpResponseStart:
	case *generated.PeerMessage_HttpResponseBody:
	case *generated.PeerMessage_HttpResponseEnd:
		break
	default:
		return nil
	}

	rawRequestID := interfaces.RequestIDOf(msg)
	requestID := httpRequestID{
		destNodeID: config.NodeID(rawRequestID.EgressNodeId),
		sequence:   uint64(rawRequestID.Sequence),
	}

	var state *inFlightHTTP
	func() {
		exec.m.Lock()
		defer exec.m.Unlock()

		state = exec.inflight[requestID]
	}()
	if state == nil {
		return fmt.Errorf("received HTTP response data but corresponding HTTP request not found (node ID: %s, seq: %v)", requestID.destNodeID, requestID.sequence)
	}
	return state.HandleResponse(ctx, msg)
}
