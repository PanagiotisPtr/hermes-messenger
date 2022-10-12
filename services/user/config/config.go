package config

import (
	"flag"
	"os"

	"github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
	"github.com/spf13/viper"
)

type ServiceConfig struct {
	Port           int
	GRPCReflection bool
}

// Application config
type Config struct {
	Service ServiceConfig
	Redis   *redis.Options
	Mongo   mongoutils.Config
}

// loadConfig loads the configuration from a file
func loadConfig(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	viper.AddConfigPath(".")

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Provides the application config
func ProvideConfig() (*Config, error) {
	configFilename := *flag.String("config", "config.dev.yml", "configuration file")
	flag.Parse()

	if os.Getenv("CONFIG_FILE") != "" {
		configFilename = os.Getenv("CONFIG_FILE")
	}

	return loadConfig(configFilename)
}

func ProvideTestConfig() (*Config, error) {
	return loadConfig("config.test.yml")
}
