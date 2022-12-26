package forwarder

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
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
	ctx := logger.Wrap(req.Context(), "http", fmt.Sprintf("%s %s %s", req.Method, req.Host, req.RequestURI))
	msgOrder := int64(0)

	logger.GetFrom(ctx).Debugw("Start forwarding HTTP request")
	dest, reqID, err := fw.startRequest(ctx, req, &msgOrder)
	if err != nil {
		return nil, err
	}
	ctx = logger.Wrap(ctx, "route", dest)
	logger.GetFrom(ctx).Debugw("Started and routed HTTP request forwarding")

	// FIXME: Implement
	return nil, fmt.Errorf("Not implemented yet")
}

func (fw *httpForwarder) startRequest(ctx context.Context, req *http.Request, msgOrder *int64) (config.NodeID, int64, error) {
	if *msgOrder != 0 {
		panic("1st message's msgOrder must be 0")
	}
	reqID := atomic.AddInt64(&fw.reqIDGen, +1)
	msg := &generated.PeerMessage{
		Message: &generated.PeerMessage_HttpRequestStart{
			HttpRequestStart: &generated.HttpRequestStart{
				Identity: &generated.HttpMessageIdentity{
					RequestId: reqID,
					MsgOrder:  *msgOrder,
				},
				Hostname: req.Host,
				Method:   req.Method,
				Path:     req.RequestURI,
			},
		},
	}
	*msgOrder++

	dest, err := fw.router.Route(ctx, msg)
	if err != nil {
		return "", reqID, err
	}

	fw.router.Deliver(ctx, fw.router.NodeID(), dest, msg)
	return dest, reqID, nil
}
