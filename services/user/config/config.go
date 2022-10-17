package config

import (
	"flag"
	"fmt"
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

	if os.Getenv("CONFIG") != "" {
		configFilename = os.Getenv("CONFIG")
	}

	return loadConfig(configFilename)
}

func ProvideTestConfig() (*Config, error) {
	configFilename := *flag.String("test-config", "config.test.yml", "configuration file")
	flag.Parse()

	if os.Getenv("TEST_CONFIG") != "" {
		configFilename = os.Getenv("TEST_CONFIG")
	}

	cfg, err := loadConfig(configFilename)

	panic(fmt.Errorf("MONGO_URI: %s", cfg.Mongo.Uri))

	return cfg, err
}
