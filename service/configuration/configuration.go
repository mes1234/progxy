package configuration

import (
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/mes1234/progxy/internal/dto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	config dto.Config
	once   sync.Once
	logger logrus.Logger
)

func init() {
	once.Do(func() {
		ReadConfig()
		SetupLogger()
	})
}

func SetupLogger() {

	logger = *logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)

	logger.SetLevel(logrus.DebugLevel)

}

func ReadConfig() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("progxyconfig.yaml")
	viper.SetConfigType("yaml")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Debug("Config file changed:", e.Name)
		port := viper.GetInt("port")
		logger.WithField("port", port).Info("Started progxy")
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&config)

	if err != nil {
		logger.Panic("Unable to read Configuration")
	}
}

type ConfigurationService interface {
	GetDestinations() *map[string]dto.Destination
	GetLogger() *logrus.Logger
}

type configurationService struct {
	configuration dto.Config
}

func (cs *configurationService) GetDestinations() *map[string]dto.Destination {
	return &cs.configuration.Destintations
}

func (cs configurationService) GetLogger() *logrus.Logger {
	return &logger
}

func NewConfigurationService() ConfigurationService {
	return &configurationService{
		configuration: config,
	}
}
