package forwarder

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type httpForwarder struct {
	router   interfaces.Router
	reqIDGen int64
}

func NewHTTPForwarder(router interfaces.Router) interfaces.HTTPForwarder {
	return &httpForwarder{router: router}
}

func (fw *httpForwarder) Close(ctx context.Context) error {
	return nil
}

func (fw *httpForwarder) RoundTrip(req *http.Request) (*http.Response, error) {
	dest, reqID, err := fw.startRequest(req)
	if err != nil {
		return nil, err
	}

	// FIXME: Implement
	return nil, fmt.Errorf("Not implemented yet")
}

func (fw *httpForwarder) startRequest(req *http.Request) (config.NodeID, *generated.RequestID, error) {
	reqID := &generated.RequestID{
		Sequence: atomic.AddInt64(&fw.reqIDGen, +1),
	}

	msg := &generated.PeerMessage{
		Message: &generated.PeerMessage_HttpRequestStart{
			HttpRequestStart: &generated.HttpRequestStart{
				RequestId: reqID,
				Hostname:  req.Host,
				Method:    req.Method,
				Path:      req.RequestURI,
			},
		},
	}

	dest, err := fw.router.Route(req.Context(), msg)
	if err != nil {
		return "", reqID, err
	}

	fw.router.Deliver(
		req.Context(),
		fw.router.NodeID(), dest,
		&generated.PeerMessage{
			Message: &generated.PeerMessage_HttpRequestStart{HttpRequestStart: msg},
		},
	)
	return dest, reqID, nil
}
