package main

import (
	"net"
	"time"

	"github.com/mes1234/progxy/internal/worker"
	"github.com/mes1234/progxy/service/configuration"
)

func main() {

	cs := configuration.NewConfigurationService()

	cs.Init(configuration.DefaultPath, configuration.NoActionCallback)

	destinations := cs.GetDestinations()

	adapters := make(map[string]worker.TcpAdapter)

	log := cs.GetLogger()

	log.Debug("Got config ", len(destinations))

	for key, destination := range destinations {
		adapters[key] = worker.NewTcpAdaper(destination.Proxied, destination.Port, net.Listen, net.Dial, net.LookupIP)
	}

	time.Sleep(1000 * time.Second)

}
