package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/mes1234/progxy/internal/worker"
	"github.com/mes1234/progxy/service/configuration"
)

func main() {

	// Create main context
	mainCtx := context.Background()

	appCtx, appCancel := context.WithCancel(mainCtx)

	// appCtx, _ = context.WithDeadline(appCtx, time.Now().Add(time.Second*10))

	// exit procedure
	exit(appCancel)

	cs := configuration.NewConfigurationService()

	cs.Init(configuration.DefaultPath, configuration.NoActionCallback)

	cs.GetLogger().Info("PROGXY bootstrap")

	run(cs, appCtx)

	<-mainCtx.Done()
}

func exit(appCancel context.CancelFunc) {
	c := make(chan os.Signal, 10)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		appCancel()
		os.Exit(0)
	}()
}

func run(cs configuration.ConfigurationService, ctx context.Context) {

	destinations := cs.GetDestinations()

	logger := cs.GetLogger()

	logger.Info("Started Progxy HAVE FUN")

	adapters := make(map[string]worker.TcpAdapter)

	for key, destination := range destinations {
		adapters[key] = worker.NewTcpAdaper(destination.Proxied, destination.Port, net.Listen, net.Dial, net.LookupIP, logger, ctx)
	}
}
