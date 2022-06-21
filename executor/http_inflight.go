package executor

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type inFlightHTTP struct {
	id   httpRequestID
	dest config.NodeID
	req  *http.Request

	// No need to lock mutex to read/write this.
	responseStarted chan (*generated.HttpResponseStart)

	responseBodyLock   sync.Mutex
	responseBodySource *io.PipeWriter
	responseBodySink   *io.PipeReader
}

func newInflightHTTP(req *http.Request, dest config.NodeID) *inFlightHTTP {
	state := &inFlightHTTP{
		req:  req,
		dest: dest,

		responseStarted: make(chan *generated.HttpResponseStart, 1),
	}
	state.responseBodySink, state.responseBodySource = io.Pipe()
	return state
}

func (state *inFlightHTTP) Close() {
	state.CloseResponseBody(io.ErrUnexpectedEOF)
}

func (state *inFlightHTTP) CloseResponseBody(err error) {
	state.responseBodyLock.Lock()
	defer state.responseBodyLock.Unlock()
	state.responseBodySource.CloseWithError(err)
}

func (state *inFlightHTTP) SendRequest(ctx context.Context, router interfaces.Router) error {
	// FIXME: Implement
	// TODO: If error, must do state.close()
	return nil
}

// AwaitResponseHeader await until the start of response body.
// Caller must read http.Response.Body to drain response body.
func (state *inFlightHTTP) AwaitResponseHeader(ctx context.Context) (*http.Response, error) {
	var responseStart *generated.HttpResponseStart
	select {
	case <-ctx.Done():
		state.Close()
		return nil, ctx.Err()
	case responseStart = <-state.responseStarted:
		break
	}

	// FIXME: Implement
	// TODO: If error, must do state.close()
	return nil, nil
}

func (state *inFlightHTTP) HandleResponse(ctx context.Context, msg interfaces.Message) error {
	switch msg := msg.Message.(type) {
	case *generated.PeerMessage_HttpResponseStart:
		state.responseStarted <- msg.HttpResponseStart
		return nil
	case *generated.PeerMessage_HttpResponseBody:
		return state.HandleResponseBody(msg.HttpResponseBody.Data)
	case *generated.PeerMessage_HttpResponseEnd:
		state.CloseResponseBody(io.EOF) // Successful EoF
		return nil
	default:
		state.Close()
		return fmt.Errorf("unexpected message received, closed in-flight HTTP request")
	}
}

func (state *inFlightHTTP) HandleResponseBody(data []byte) error {
	// TODO: Implement some timeout (response body consumer may not drain the body properly)
	// TODO: Implement some buffering to not block source
	state.responseBodyLock.Lock()
	defer state.responseBodyLock.Unlock()
	_, err := state.responseBodySource.Write(data)
	return err
}
