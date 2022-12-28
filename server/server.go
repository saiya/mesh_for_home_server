package server

import (
	"context"
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/egress"
	"github.com/saiya/mesh_for_home_server/ingress"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering"
	"github.com/saiya/mesh_for_home_server/router"
)

type Server struct {
	ctx      context.Context
	ctxClose context.CancelFunc

	//
	// Fields ordered by initialization/close order
	//

	router         interfaces.Router
	egressHandlers []interfaces.MessageHandler
	peeringServers []interfaces.PeeringServer
	peeringClients []interfaces.PeeringClient
	forwarders     []interfaces.Forwarder
	ingresses      []interfaces.Ingress
}

func StartServer(config *config.ServerConfig) (*Server, error) {
	router := router.NewRouter(config.Hostname)

	ctx, ctxClose := context.WithCancel(context.Background())

	advFn, egressHandlers, err := egress.StartEgress(config.Egress, router)
	if err != nil {
		ctxClose()
		return nil, fmt.Errorf("failed to initialize egress: %w", err)
	}
	router.SetAdvertisementProvider(advFn)

	peeringServers, peeringClients, err := peering.StartPeering(ctx, config.Perring, router)
	if err != nil {
		ctxClose()
		return nil, fmt.Errorf("failed to initialize ingress: %w", err)
	}

	ingresses, forwarders, err := ingress.StartIngress(config.Ingress, router)
	if err != nil {
		ctxClose()
		return nil, fmt.Errorf("failed to initialize ingress: %w", err)
	}

	return &Server{
		ctx: ctx, ctxClose: ctxClose,

		router:         router,
		egressHandlers: egressHandlers,
		peeringServers: peeringServers,
		peeringClients: peeringClients,
		forwarders:     forwarders,
		ingresses:      ingresses,
	}, nil
}

func (srv *Server) Close(ctx context.Context) {
	catch := func(err error) {
		if err != nil {
			logger.GetFrom(ctx).Info("Failed to close: "+err.Error(), "err", err)
		}
	}

	defer func() { srv.ctxClose() }()

	for _, ingress := range srv.ingresses {
		catch(ingress.Close(ctx))
	}
	for _, forwarder := range srv.forwarders {
		catch(forwarder.Close(ctx))
	}
	catch(srv.router.Close(ctx))
	for _, peeringClient := range srv.peeringClients {
		catch(peeringClient.Close(ctx))
	}
	for _, peeringServer := range srv.peeringServers {
		catch(peeringServer.Close(ctx))
	}
	for _, egressHandler := range srv.egressHandlers {
		catch(egressHandler.Close(ctx))
	}
}

func (srv *Server) AwaitTerminate() {
	<-srv.ctx.Done()
}
