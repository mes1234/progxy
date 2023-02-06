package main

import (
	"os"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("progxyconfig.yaml")
	viper.SetConfigType("yaml")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Debug("Config file changed:", e.Name)
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	port := viper.GetInt("port")

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)

	log.WithField("port", port).Info("Started progxy")
}
