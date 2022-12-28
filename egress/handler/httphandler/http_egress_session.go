package httphandler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/messagewindow"
	"github.com/saiya/mesh_for_home_server/peering/proto"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type httpEgressSessionID struct {
	peer      config.NodeID
	requestID int64
}

type httpEgressSession struct {
	httpHandler     *httpHandler
	egress          *httpEgress
	responseTimeout *config.HTTPTimeout

	ID httpEgressSessionID

	ctx       context.Context
	ctxCancel context.CancelFunc
	req       *http.Request

	bodyWindow messagewindow.MessageWindow[int64, interface{}] // []byte or err

	replyMsgOrder int64
}

const httpReqBodyChunkSize = 512 * 1024

func newHTTPEgressSession(ctx context.Context, h *httpEgress, req *generated.HttpRequestStart, id httpEgressSessionID) (*httpEgressSession, error) {
	url, err := url.Parse(fmt.Sprintf("%s%s", h.server, req.Path))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL to forward HTTP request: %v", err)
	}

	ctx, ctxCancel := context.WithCancel(logger.Wrap(ctx, "url", url, "hostname", req.Hostname, "request-id", id.requestID))
	sess := &httpEgressSession{
		httpHandler:     h.httpHandler,
		egress:          h,
		responseTimeout: h.responseTimeout,

		ID: id,

		ctx: ctx, ctxCancel: ctxCancel,
		req: (&http.Request{
			Method: req.Method,
			URL:    url,
			Host:   req.Hostname,
			Header: proto.ToHTTPHeaders(req.Headers),
		}).WithContext(ctx),

		bodyWindow: messagewindow.NewMessageWindow[int64, interface{}](),
	}
	sess.req.Body = &httpEgressSessionReqBody{sess: sess}
	return sess, nil
}

func (sess *httpEgressSession) Close(ctx context.Context) error {
	sess.httpHandler.forgetSession(sess.ID)
	sess.bodyWindow.Close()
	return nil
}

func (sess *httpEgressSession) start() {
	go func() {
		logger.GetFrom(sess.ctx).Debugw("fowarding egress HTTP request...")
		res, err := sess.egress.httpClient.Do(sess.req)
		if res != nil && res.Body != nil { // Even on error, should close response body
			defer func() {
				err := res.Body.Close()
				if err != nil {
					logger.GetFrom(sess.ctx).Debugf("failed to close response body: %w", err)
				}
			}()
		}
		if err != nil {
			sess.handleRequestFailure(err)
			return
		}
		sess.forwardResponse(res)
	}()
}

func (sess *httpEgressSession) handle(ctx context.Context, msg *generated.HttpMessage) error {
	logger.GetFrom(ctx).Debugw("Handling HTTP egress message")
	if body := msg.GetHttpRequestBody(); body != nil {
		return sess.bodyWindow.Send(msg.Identity.MsgOrder-1, body.Data)
	} else if end := msg.GetHttpRequestEnd(); end != nil {
		sess.bodyWindow.Close()
		return nil
	} else if end := msg.GetHttpRequestAbnormalEnd(); end != nil {
		err := sess.bodyWindow.Send(msg.Identity.MsgOrder-1, fmt.Errorf("remote peer terminated HTTP request"))
		sess.bodyWindow.Close() // Anyway we should close window to notify end of request
		return err
	} else {
		return fmt.Errorf("HTTP egress forwarder encountered unexpected message: %v", msg.Message)
	}
}

type httpEgressSessionReqBody struct {
	sess *httpEgressSession

	bufferedBody []byte
}

func (b *httpEgressSessionReqBody) Close() error { // io.Closer
	// This Close() means end of the client request read, NOT end of the HTTP request-response lifecycle
	b.sess.bodyWindow.Close()
	return nil
}

func (b *httpEgressSessionReqBody) Read(p []byte) (n int, err error) { // io.Reader
	if len(b.bufferedBody) > 0 {
		copied := copy(p, b.bufferedBody)
		b.bufferedBody = b.bufferedBody[copied:]
		return copied, nil
	}

	ctx, ctxCancel := context.WithTimeout(b.sess.ctx, b.sess.responseTimeout.BodyTimeout())
	defer ctxCancel()

	_, msg, err := b.sess.bodyWindow.Consume(ctx)
	if err != nil {
		return 0, err
	}
	if body, ok := msg.([]byte); ok {
		copied := copy(p, body)
		b.bufferedBody = body[copied:]
		return copied, nil
	}
	return 0, msg.(error)
}

func (sess *httpEgressSession) reply(msg *generated.HttpMessage) {
	msg.Identity = &generated.HttpMessageIdentity{
		RequestId: sess.ID.requestID,
		MsgOrder:  atomic.AddInt64(&sess.replyMsgOrder, 1) - 1,
	}

	router := sess.httpHandler.router
	router.Deliver(sess.ctx, router.NodeID(), sess.ID.peer, &generated.PeerMessage{Message: &generated.PeerMessage_Http{Http: msg}})
}

func (sess *httpEgressSession) forwardResponse(res *http.Response) {
	logger.GetFrom(sess.ctx).Debugw("Got HTTP response from forwarding target", "status", res.Status)

	sess.reply(&generated.HttpMessage{
		Message: &generated.HttpMessage_HttpResponseStart{
			HttpResponseStart: &generated.HttpResponseStart{
				Status: res.Status, StatusCode: int32(res.StatusCode),
				Proto: res.Proto, ProtoMajor: int32(res.ProtoMajor), ProtoMinor: int32(res.ProtoMinor),

				ContentLength:     res.ContentLength,
				TransferEncodings: res.TransferEncoding,
				Headers:           proto.FromHTTPHeaders(res.Header),
			},
		},
	})

	for {
		buf := make([]byte, httpReqBodyChunkSize) // Must avoid reusing buffer, otherwise succeeding message's content will be polluted
		bytes, err := res.Body.Read(buf)
		if bytes > 0 {
			sess.reply(&generated.HttpMessage{
				Message: &generated.HttpMessage_HttpResponseBody{
					HttpResponseBody: &generated.HttpResponseBody{
						Data: buf[:bytes],
					},
				},
			})
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				sess.reply(&generated.HttpMessage{
					Message: &generated.HttpMessage_HttpResponseEnd{
						HttpResponseEnd: &generated.HttpResponseEnd{
							Trailers: proto.FromHTTPHeaders(res.Trailer),
						},
					},
				})
			} else {
				logger.GetFrom(sess.ctx).Warnw("Failed to read response from forwarding target: "+err.Error(), "err", err)
				sess.reply(&generated.HttpMessage{
					Message: &generated.HttpMessage_HttpResponseAbnormalEnd{
						HttpResponseAbnormalEnd: &generated.HttpResponseAbnormalEnd{},
					},
				})
			}
			break
		}
	}
}

func (sess *httpEgressSession) handleRequestFailure(err error) {
	logger.GetFrom(sess.ctx).Warnw("Failed to forward HTTP request: "+err.Error(), "err", err)

	// Send 502 response
	responseBody := []byte("Failed to forward HTTP request.")
	sess.reply(&generated.HttpMessage{
		Message: &generated.HttpMessage_HttpResponseStart{
			HttpResponseStart: &generated.HttpResponseStart{
				Status: "502 Bad Gateway", StatusCode: 502,
				Proto: "HTTP/1.1", ProtoMajor: 1, ProtoMinor: 1,

				ContentLength: int64(len(responseBody)),
				Headers: proto.FromHTTPHeaders(map[string][]string{
					"Content-Type": {"text/plain"},
				}),
			},
		},
	})
	sess.reply(&generated.HttpMessage{
		Message: &generated.HttpMessage_HttpResponseBody{
			HttpResponseBody: &generated.HttpResponseBody{
				Data: responseBody,
			},
		},
	})
	sess.reply(&generated.HttpMessage{
		Message: &generated.HttpMessage_HttpResponseEnd{
			HttpResponseEnd: &generated.HttpResponseEnd{},
		},
	})
}
