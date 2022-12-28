package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/server"
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
	if err := cfgInput.Close(); err != nil {
		logger.Get().Warnf("cannot close configuration file: %w", err)
	}

	srv, err := server.StartServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize server: %w", err)
	}

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	go func() {
		<-interruptCh
		srv.Close(context.Background())
	}()

	srv.AwaitTerminate()
	return nil
}
