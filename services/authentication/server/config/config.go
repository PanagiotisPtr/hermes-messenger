package config

import (
	"flag"
	"path/filepath"
	"time"

	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/spf13/viper"
)

// The application config
type Config struct {
	ListenPort           int           `mapstructure:"listenPort"`
	GRPCReflection       bool          `mapstructure:"grpcReflection"`
	RedisAddress         string        `mapstructure:"redisAddress"`
	RedisPassword        string        `mapstructure:"redisPassword"`
	RedisDatabase        int           `mapstructure:"redisDatabase"`
	AMQPURI              string        `mapstructure:"amqpURI"`
	RefreshTokenDuration time.Duration `mapstructure:"refreshTokenDuration"`
	AccessTokenDuration  time.Duration `mapstructure:"accessTokenDuration"`
}

// Provides the application config
func ProvideConfig() (*Config, error) {
	configPath := flag.Lookup("config").Value.(flag.Getter).Get().(string)
	viper.SetConfigName(filepath.Base(configPath))
	viper.AddConfigPath(filepath.Dir(configPath))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	// we get these values from secrets so they shouldn't be in the config
	cfg.RedisAddress = utils.GetEnvVariableString("REDIS_ADDRESS", "localhost:6379")
	cfg.RedisPassword = utils.GetEnvVariableString("REDIS_PASSWORD", "")
	cfg.RedisDatabase = utils.GetEnvVariableInt("REDIS_DB", 0)
	cfg.AMQPURI = utils.GetEnvVariableString("AMQP_URI", "amqp://guest:guest@localhost:5672")

	return cfg, nil
}
