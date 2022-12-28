package httpingress

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/tlshelper"
)

type httpIngress struct {
	config *config.HTTPIngressConfig

	port    int
	hserver *http.Server

	closeWait *sync.WaitGroup
	closeErr  error
}

func NewHTTPIngress(
	config *config.HTTPIngressConfig,
	// Default handler that should handle requests.
	// In this system, this will be wrapper of httputil.ReverseProxy.
	defaultHandler func(http.ResponseWriter, *http.Request) error,
) (interfaces.Ingress, error) {
	var listenAddr string
	if config.Listen != "" {
		listenAddr = config.Listen
	} else {
		listenAddr = ":0"
	}

	tcpListener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen \"%s\": %w", listenAddr, err)
	}

	h := httpIngress{
		config:    config,
		port:      tcpListener.Addr().(*net.TCPAddr).Port,
		closeWait: &sync.WaitGroup{},
		hserver: &http.Server{
			Addr:    listenAddr,
			Handler: httpHandler(config, defaultHandler),

			ReadTimeout:       config.RequestTimeout.BodyTimeout(),
			ReadHeaderTimeout: config.RequestTimeout.HeaderTimeout(),
			WriteTimeout:      config.ResponseTimeout.BodyTimeout(),
			IdleTimeout:       config.GetKeepAliveTimeout(),
		},
	}
	if config.TLS != nil {
		tlsConfig, err := tlshelper.BuildTLSServerConfig(config.TLS)
		if err != nil {
			return nil, fmt.Errorf("failed to setup TLS: %w", err)
		}
		h.hserver.TLSConfig = tlsConfig
	}

	h.closeWait.Add(1)
	go func() {
		defer h.closeWait.Done()

		var err error
		if h.hserver.TLSConfig != nil {
			logger.Get().Infow("Starting HTTPS ingress", "port", h.port)
			err = h.hserver.ServeTLS(tcpListener, "", "")
		} else {
			logger.Get().Infow("Starting HTTP ingress", "port", h.port)
			err = h.hserver.Serve(tcpListener)
		}
		if err != http.ErrServerClosed {
			logger.Get().Errorw("Failed to start HTTP(S) ingress: "+err.Error(), "err", err)
			h.closeErr = err
		} else {
			logger.Get().Infow("HTTP(S) ingress closed", "err", err)
		}
	}()
	return &h, nil
}

func (h *httpIngress) Close(ctx context.Context) error {
	warnIfError("Shutdown() returned error", h.hserver.Shutdown(ctx))
	h.closeWait.Wait()
	return h.closeErr
}

func (h *httpIngress) String() string {
	var tlsDesc string
	if h.config.TLS == nil {
		tlsDesc = "none"
	} else {
		if h.config.TLS.ClientCertCAFile == "" {
			tlsDesc = "enabled"
		} else {
			tlsDesc = "require_client_cert"
		}
	}

	return fmt.Sprintf(
		"HTTPIngress{ port: %d, TLS: %s }",
		h.port,
		tlsDesc,
	)
}
