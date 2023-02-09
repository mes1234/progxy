package main

import (
	"time"

	"github.com/mes1234/progxy/service/configuration"
)

func main() {

	cs := configuration.NewConfigurationService()

	cs.Init(configuration.DefaultPath, configuration.NoActionCallback)

	config := cs.GetDestinations()

	log := cs.GetLogger()

	log.Debug("Got config ", len(*config))

	time.Sleep(1000 * time.Second)

}
