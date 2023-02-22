package main

import (
	"context"
	"net"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/mes1234/progxy/internal/worker"
	"github.com/mes1234/progxy/service/configuration"
)

func main() {

	mainCtx := context.Background()

	cs := configuration.NewConfigurationService()

	runCtx, cancelFunc := context.WithCancel(mainCtx)

	cs.Init(configuration.DefaultPath, reloadFuncGenerator(cs, cancelFunc, mainCtx))

	logger := cs.GetLogger()

	logger.Info("PROGXY bootstrap")

	run(cs, runCtx, func() {})

	<-mainCtx.Done()

}

func run(cs configuration.ConfigurationService, ctx context.Context, cancelPrior context.CancelFunc) {

	cancelPrior()

	time.Sleep(1 * time.Second)

	destinations := cs.GetDestinations()

	logger := cs.GetLogger()

	logger.Debug("Got config ", len(destinations))

	logger.Info("Started Progxy HAVE FUN")

	adapters := make(map[string]worker.TcpAdapter)

	for key, destination := range destinations {
		adapters[key] = worker.NewTcpAdaper(destination.Proxied, destination.Port, net.Listen, net.Dial, net.LookupIP, logger, ctx)
	}
}

func reloadFuncGenerator(cs configuration.ConfigurationService, cancelPrior context.CancelFunc, mainCtx context.Context) func(e fsnotify.Event) {

	return func(e fsnotify.Event) {

		runCtx, cancelFunc := context.WithCancel(mainCtx)

		cs.Init(configuration.DefaultPath, reloadFuncGenerator(cs, cancelFunc, mainCtx))

		logger := cs.GetLogger()

		logger.Info("PROGXY bootstrap")

		run(cs, runCtx, func() {})

	}
}
