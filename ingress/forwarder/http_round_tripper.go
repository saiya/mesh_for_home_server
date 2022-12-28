package forwarder

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type httpRoundTripper struct {
	fw     *httpForwarder
	router interfaces.Router

	headerTimeout time.Duration
	bodyTimeout   time.Duration
}

func (hrt *httpRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := logger.Wrap(req.Context(), "http", fmt.Sprintf("%s %s %s", req.Method, req.Host, req.RequestURI))

	// Don't close this in the end of this method, close this in the end of HTTP response body.
	fwc := hrt.fw.newSession(hrt.headerTimeout, hrt.bodyTimeout)

	logger.GetFrom(ctx).Debugw("Start forwarding HTTP request")
	if err := hrt.startRequest(ctx, fwc, req); err != nil {
		fwc.Close()
		return nil, err
	}
	ctx = logger.Wrap(ctx, "route", fwc.dest)
	logger.GetFrom(ctx).Debugw("Started and routed HTTP request forwarding")

	err := hrt.forwardReqBody(ctx, fwc, req)
	hrt.endOfReqBody(ctx, fwc, err) // Always send end of request, even though it terminated
	if err != nil {
		fwc.Close()
		return nil, err
	}
	logger.GetFrom(ctx).Debugw("Forwarded HTTP request, waiting response...")

	return hrt.readResponse(ctx, fwc)
}

func (hrt *httpRoundTripper) startRequest(ctx context.Context, fwc *httpForwardingSession, req *http.Request) error {
	innerMsg := &generated.HttpRequestStart{
		Hostname: req.Host,
		Method:   req.Method,
		Path:     req.RequestURI,
		Headers:  make([]*generated.HttpHeader, 0, len(req.Header)),
	}
	for k, v := range req.Header {
		innerMsg.Headers = append(innerMsg.Headers, &generated.HttpHeader{Name: k, Values: v})
	}
	msg := &generated.PeerMessage{
		Message: &generated.PeerMessage_Http{
			Http: &generated.HttpMessage{
				Identity: fwc.NextMsgID(),
				Message: &generated.HttpMessage_HttpRequestStart{
					HttpRequestStart: innerMsg,
				},
			},
		},
	}

	dest, err := hrt.router.Route(ctx, msg)
	if err != nil {
		return err
	}
	fwc.dest = dest
	fwc.startListener()

	hrt.router.Deliver(ctx, fwc.from, dest, msg)
	return nil
}

func (hrt *httpRoundTripper) forwardReqBody(ctx context.Context, fwc *httpForwardingSession, req *http.Request) error {
	if req.Body == nil {
		return nil
	}

	buf := make([]byte, httpReqBodyChunkSize)
	for {
		n, err := req.Body.Read(buf)
		if n > 0 {
			hrt.router.Deliver(ctx, fwc.from, fwc.dest, &generated.PeerMessage{
				Message: &generated.PeerMessage_Http{
					Http: &generated.HttpMessage{
						Identity: fwc.NextMsgID(),
						Message: &generated.HttpMessage_HttpRequestBody{
							HttpRequestBody: &generated.HttpRequestBody{
								Data: buf[0:n],
							},
						},
					},
				},
			})
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			} else {
				return fmt.Errorf("error while forwarding HTTP request body to mesh: %w", err)
			}
		}
	}
	return nil
}

func (hrt *httpRoundTripper) endOfReqBody(ctx context.Context, fwc *httpForwardingSession, err error) {
	if err != nil {
		hrt.router.Deliver(ctx, fwc.from, fwc.dest, &generated.PeerMessage{
			Message: &generated.PeerMessage_Http{
				Http: &generated.HttpMessage{
					Identity: fwc.NextMsgID(),
					Message: &generated.HttpMessage_HttpRequestAbnormalEnd{
						HttpRequestAbnormalEnd: &generated.HttpRequestAbnormalEnd{},
					},
				},
			},
		})
	} else {
		hrt.router.Deliver(ctx, fwc.from, fwc.dest, &generated.PeerMessage{
			Message: &generated.PeerMessage_Http{
				Http: &generated.HttpMessage{
					Identity: fwc.NextMsgID(),
					Message: &generated.HttpMessage_HttpRequestEnd{
						HttpRequestEnd: &generated.HttpRequestEnd{},
					},
				},
			},
		})
	}
}

func (hrt *httpRoundTripper) readResponse(ctx context.Context, fwc *httpForwardingSession) (*http.Response, error) {
	res, err := hrt.readResHeader(ctx, fwc)
	if err != nil {
		fwc.Close()
		return res, err
	}

	res.Body = &httpResponseBodyReader{hrt: hrt, ctx: ctx, fwc: fwc, res: res}

	logger.GetFrom(ctx).Debugw("HTTP response started", "proto", res.Proto, "status", res.Status, "content-length", res.ContentLength)
	return res, nil
}

func (hrt *httpRoundTripper) readResHeader(ctx context.Context, fwc *httpForwardingSession) (*http.Response, error) {
	consumeCtx, consumeCancel := context.WithTimeout(ctx, hrt.headerTimeout)
	_, msg, err := fwc.msgWindow.Consume(consumeCtx)
	consumeCancel()
	if err != nil {
		return nil, fmt.Errorf("failed to await/read HTTP response header %w", err)
	}

	rs := msg.GetHttpResponseStart()
	if rs == nil {
		return nil, fmt.Errorf("Got unexpected message from peer: %v", msg.Message)
	}
	return proto.ToHTTPResponse(rs), nil
}

type httpResponseBodyReader struct {
	hrt *httpRoundTripper
	ctx context.Context
	fwc *httpForwardingSession
	res *http.Response

	bufferedBody []byte
}

func (this *httpResponseBodyReader) Close() error {
	this.fwc.Close()
	return nil
}

func (this *httpResponseBodyReader) Read(p []byte) (n int, err error) {
	if len(this.bufferedBody) > 0 {
		copied := copy(p, this.bufferedBody)
		this.bufferedBody = this.bufferedBody[copied:]
		return copied, nil
	}

	ctx := this.ctx

	consumeCtx, consumeCancel := context.WithTimeout(ctx, this.hrt.bodyTimeout)
	_, msg, err := this.fwc.msgWindow.Consume(consumeCtx)
	consumeCancel()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, fmt.Errorf("Unexpected end-of-message-sequence (missing HttpResponseEnd message) from peer")
		} else {
			return 0, fmt.Errorf("Failed to await/read HTTP response body or end from peer: %v", err)
		}
	}

	if body := msg.GetHttpResponseBody(); body != nil {
		copied := copy(p, body.Data)
		this.bufferedBody = body.Data[copied:]
		return copied, nil
	} else if end := msg.GetHttpResponseEnd(); end != nil {
		defer this.fwc.Close()

		this.res.Trailer = proto.ToHTTPHeaders(end.Trailers)
		return 0, io.EOF
	} else if end := msg.GetHttpResponseAbnormalEnd(); end != nil {
		defer this.fwc.Close()
		return 0, fmt.Errorf("Unexpected response termination from remote peer")
	} else {
		return 0, fmt.Errorf("Unexpected message from peer while reading HTTP response: %v", msg.Message)
	}
}
