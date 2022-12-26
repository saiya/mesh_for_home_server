package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/egress"
	"github.com/saiya/mesh_for_home_server/ingress"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering"
	"github.com/saiya/mesh_for_home_server/router"
)

var debugFlag = flag.Bool("debug", false, "shows DEBUG logs")

func main() {
	err := mainImpl()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func mainImpl() error {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		return fmt.Errorf("must give configuration file path in command line argument")
	}

	if *debugFlag {
		logger.EnableDebugLog()
	}

	cfgInput, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("cannot open configuration file \"%s\": %w", args[0], err)
	}
	cfg, err := config.ParseConfig(cfgInput)
	if err != nil {
		return fmt.Errorf("invalid configuration file: %w", err)
	}
	cfgInput.Close()

	srv, err := StartServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	go func() {
		<-interruptCh
		srv.Close(context.Background())
	}()
	<-srv.ctx.Done()
	return nil
}

type server struct {
	ctx      context.Context
	ctxClose context.CancelFunc

	router         interfaces.Router
	peeringServers []interfaces.PeeringServer
	peeringClients []interfaces.PeeringClient
	egressHandlers []interfaces.MessageHandler
	ingresses      []interfaces.Ingress
	forwarders     []interfaces.Forwarder
}

func StartServer(config *config.ServerConfig) (*server, error) {
	router := router.NewRouter("")

	ctx, ctxClose := context.WithCancel(context.Background())
	peeringServers, peeringClients, err := peering.StartPeering(ctx, config.Perring, router)
	if err != nil {
		ctxClose()
		return nil, fmt.Errorf("failed to initialize ingress: %w", err)
	}

	advFn, egressHandlers, err := egress.StartEgress(config.Egress, router)
	if err != nil {
		ctxClose()
		return nil, fmt.Errorf("failed to initialize egress: %w", err)
	}
	router.SetAdvertisementProvider(advFn)

	ingresses, forwarders, err := ingress.StartIngress(config.Ingress, router)
	if err != nil {
		ctxClose()
		return nil, fmt.Errorf("failed to initialize ingress: %w", err)
	}

	return &server{
		ctx: ctx, ctxClose: ctxClose,

		router:         router,
		peeringServers: peeringServers,
		peeringClients: peeringClients,
		egressHandlers: egressHandlers,
		ingresses:      ingresses,
		forwarders:     forwarders,
	}, nil
}

func (srv *server) Close(ctx context.Context) {
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
	for _, egressHandler := range srv.egressHandlers {
		catch(egressHandler.Close(ctx))
	}
	for _, peeringClient := range srv.peeringClients {
		catch(peeringClient.Close(ctx))
	}
	for _, peeringServer := range srv.peeringServers {
		catch(peeringServer.Close(ctx))
	}
}
