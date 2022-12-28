package httphandler

import (
	"context"
	"fmt"
	"sync"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
	"golang.org/x/sync/errgroup"
)

type HttpHandler interface {
	interfaces.MessageHandler
	AddEgress(c *config.HTTPEgressConfig) error
	Advertise() *generated.HttpAdvertisement
}
type httpHandler struct {
	router            interfaces.Router
	closeRouterListen func()

	egressesLock sync.RWMutex
	egresses     map[string]*httpEgress

	sessionsLock sync.RWMutex
	sessions     map[httpEgressSessionID]*httpEgressSession
}

func NewHttpHandler(router interfaces.Router) HttpHandler {
	h := &httpHandler{
		router:   router,
		egresses: make(map[string]*httpEgress),
		sessions: make(map[httpEgressSessionID]*httpEgressSession),
	}
	h.closeRouterListen = router.Listen(func(ctx context.Context, from config.NodeID, msg interfaces.Message) error {
		http := msg.GetHttp()
		if http == nil {
			return nil
		}
		sessID := httpEgressSessionID{from, http.Identity.RequestId}
		ctx = logger.Wrap(ctx, "peer", sessID.peer, "request-id", sessID.requestId)

		if req := http.GetHttpRequestStart(); req != nil {
			egress := h.findEgress(req)
			if egress == nil {
				return fmt.Errorf("no egress found to handle incoming request (host: %s)", req.Hostname)
			}

			sess, err := newHttpEgressSession(ctx, egress, req, sessID)
			if err != nil {
				return fmt.Errorf("failed to start HTTP egress request (host: %s): %v", req.Hostname, err)
			}

			// Immediately register the session, to properly handle suceeding messages
			h.sessionsLock.Lock()
			h.sessions[sessID] = sess
			h.sessionsLock.Unlock()
			logger.GetFrom(ctx).Debugw("Registered HTTP egress session")

			sess.start()
			return nil
		} else {
			h.sessionsLock.RLock()
			sess := h.sessions[sessID]
			defer h.sessionsLock.RUnlock()

			if sess == nil {
				// Because other message (e.g. HTTP ingress's message) can come, this case can happen even in normal case
				return nil
			} else {
				return sess.handle(ctx, http)
			}
		}
	})
	return h
}

func (h *httpHandler) Close(ctx context.Context) error {
	h.closeRouterListen()

	var eg errgroup.Group

	func() {
		h.sessionsLock.Lock()
		defer h.sessionsLock.Unlock()
		for key, sess := range h.sessions {
			eg.Go(func() error {
				return sess.Close(ctx)
			})
			delete(h.sessions, key)
		}
	}()

	func() {
		h.egressesLock.Lock()
		defer h.egressesLock.Unlock()
		for key, egress := range h.egresses {
			eg.Go(func() error {
				return egress.Close(ctx)
			})
			delete(h.egresses, key)
		}
	}()

	return eg.Wait()
}

func (h *httpHandler) Advertise() *generated.HttpAdvertisement {
	h.egressesLock.RLock()
	defer h.egressesLock.RUnlock()

	hostnames := make([]string, 0, len(h.egresses))
	for k := range h.egresses {
		hostnames = append(hostnames, k)
	}

	return &generated.HttpAdvertisement{
		HostnameMatchers: hostnames,
	}
}

func (h *httpHandler) AddEgress(c *config.HTTPEgressConfig) error {
	h.egressesLock.Lock()
	defer h.egressesLock.Unlock()

	if h.egresses[c.Host] != nil {
		return fmt.Errorf("HTTP egress for hostname pattern \"%s\" already exists. Could not register duplicated egress", c.Host)
	}

	egress := newHttpEgress(h, c)
	h.egresses[c.Host] = egress
	return nil
}

func (h *httpHandler) findEgress(req *generated.HttpRequestStart) *httpEgress {
	h.egressesLock.RLock()
	defer h.egressesLock.RUnlock()

	var egress *httpEgress
	egressPriority := int64(-1)
	for _, eg := range h.egresses {
		p := eg.hostMatcher(req.Hostname)
		if p > egressPriority {
			egress = eg
			egressPriority = p
		}
	}
	return egress
}

// forgetSession won't close session itself, just remove it from look up table
func (h *httpHandler) forgetSession(sessID httpEgressSessionID) {
	logger.Get().Debugw("Removing HTTP egress session", "peer", sessID.peer, "request-id", sessID.requestId)

	h.sessionsLock.Lock()
	defer h.sessionsLock.Unlock()

	delete(h.sessions, sessID)
}
