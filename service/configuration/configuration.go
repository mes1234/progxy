package configuration

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/mes1234/progxy/internal/dto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func (cs *configurationService) setupLogger() {

	cs.logger = *logrus.New()

	cs.logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	cs.logger.SetOutput(os.Stdout)

	cs.logger.SetLevel(logrus.DebugLevel)

}

type RealoadCallback func(e fsnotify.Event)

func NoActionCallback(e fsnotify.Event) {}

var DefaultPath = "."

func (cs *configurationService) Init(path string, realoadCallback RealoadCallback) {

	cs.readConfig(path, realoadCallback)
	cs.setupLogger()

}

func (cs *configurationService) readConfig(path string, realoadCallback RealoadCallback) {

	if path != "" {
		viper.AddConfigPath(path)
	} else {
		viper.AddConfigPath(".")
	}

	viper.SetConfigFile("progxyconfig.yaml")
	viper.SetConfigType("yaml")
	viper.WatchConfig()
	viper.OnConfigChange(realoadCallback)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&cs.config)

	if err != nil {
		cs.logger.Panic("Unable to read Configuration")
	}
}

type ConfigurationService interface {
	GetDestinations() map[string]dto.Destination
	GetLogger() *logrus.Logger
	Init(path string, realoadCallback RealoadCallback)
}

type configurationService struct {
	logger logrus.Logger
	config dto.Config
}

func (cs *configurationService) GetDestinations() map[string]dto.Destination {
	return cs.config.Destintations
}

func (cs *configurationService) GetLogger() *logrus.Logger {
	return &cs.logger
}

func NewConfigurationService() ConfigurationService {
	return &configurationService{}
}
