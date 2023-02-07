package main

import (
	"time"

	"github.com/mes1234/progxy/service/configuration"
)

func main() {

	config := configuration.NewConfigurationService().GetDestinations()

	log := configuration.NewConfigurationService().GetLogger()

	log.Debug("Got config ", len(*config))

	time.Sleep(1000 * time.Second)

}
