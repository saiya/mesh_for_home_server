package forwarder

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/peering/messagewindow"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
	"golang.org/x/exp/maps"
)

type httpForwarder struct {
	router            interfaces.Router
	closeRouterListen func()

	reqIDGen int64 // This ID must be unique within httpForwarder, not httpForwardingRountTripper

	listenersLock sync.Mutex
	listeners     map[httpForwarderListnerID]*httpForwardingSession
}

type httpForwarderListnerID struct {
	peer  config.NodeID
	reqID int64
}

const httpReqBodyChunkSize = 512 * 1024

func NewHTTPForwarder(router interfaces.Router) interfaces.HTTPForwarder {
	fw := &httpForwarder{
		router:    router,
		listeners: make(map[httpForwarderListnerID]*httpForwardingSession),
	}
	fw.closeRouterListen = router.Listen(func(ctx context.Context, from config.NodeID, msg interfaces.Message) error {
		if http := msg.GetHttp(); http != nil {
			listenerID := httpForwarderListnerID{from, http.Identity.RequestId}

			fw.listenersLock.Lock()
			listener := fw.listeners[listenerID]
			fw.listenersLock.Unlock()

			if listener == nil {
				return fmt.Errorf("HTTP forwarder cannot find the in-flight request: %v", http.Identity.RequestId)
			} else {
				return listener.handle(http)
			}
		}
		return nil
	})
	return fw
}

func (fw *httpForwarder) Close(ctx context.Context) error {
	fw.closeRouterListen()

	fw.listenersLock.Lock()
	listeners := maps.Values(fw.listeners)
	fw.listenersLock.Unlock()

	for _, listener := range listeners {
		listener.Close()
	}
	return nil
}

func (fw *httpForwarder) newSession(headerTimeout, bodyTimeout time.Duration) *httpForwardingSession {
	fwc := &httpForwardingSession{
		fw:            fw,
		reqID:         atomic.AddInt64(&fw.reqIDGen, +1),
		msgOrder:      0,
		msgWindow:     messagewindow.NewMessageWindow[int64, *generated.HttpMessage](),
		from:          fw.router.NodeID(),
		headerTimeout: headerTimeout,
		bodyTimeout:   bodyTimeout,
	}
	fwc.startListener()
	return fwc
}

func (fw *httpForwarder) NewRoundTripper(cfg *config.HTTPIngressConfig) http.RoundTripper {
	return &httpRoundTripper{
		fw, fw.router,
		cfg.ResponseTimeout.HeaderTimeout(), cfg.ResponseTimeout.BodyTimeout(),
	}
}

func (fw *httpForwarder) addListener(fwc *httpForwardingSession) {
	fw.listenersLock.Lock()
	defer fw.listenersLock.Unlock()

	fw.listeners[httpForwarderListnerID{fwc.dest, fwc.reqID}] = fwc
}

func (fw *httpForwarder) removeListener(fwc *httpForwardingSession) {
	fw.listenersLock.Lock()
	defer fw.listenersLock.Unlock()

	delete(fw.listeners, httpForwarderListnerID{fwc.dest, fwc.reqID})
}

type httpForwardingSession struct {
	fw *httpForwarder

	reqID    int64
	msgOrder int64

	msgWindow messagewindow.MessageWindow[int64, *generated.HttpMessage]

	from config.NodeID
	dest config.NodeID

	headerTimeout time.Duration
	bodyTimeout   time.Duration
}

// Close completely this session.
// This method make sure all resources of this session to be freed.
func (fwc *httpForwardingSession) Close() {
	fwc.msgWindow.Close()
	fwc.fw.removeListener(fwc)
}

func (fwc *httpForwardingSession) startListener() {
	fwc.fw.addListener(fwc)
	// TODO: Run session reaper, to forcibly close abandoned session
}

func (fwc *httpForwardingSession) handle(msg *generated.HttpMessage) error {
	if err := fwc.msgWindow.Send(msg.Identity.MsgOrder, msg); err != nil {
		return fmt.Errorf("failed to handle HTTP message from peer: %v", err)
	}
	return nil
}

// NextMsgID issues message identity. Returns unique & sequential identity for each call.
func (fwc *httpForwardingSession) NextMsgID() *generated.HttpMessageIdentity {
	id := &generated.HttpMessageIdentity{
		RequestId: fwc.reqID,
		MsgOrder:  fwc.msgOrder,
	}
	fwc.msgOrder++
	return id
}
