package redisutils

import (
	"strings"

	redis "github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

// ProvideRedisClient provides redis client
func ProvideRedisClient(cfg *RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}

func ProvideRedisConfig(cl *utils.ConfigLocation) (*RedisConfig, error) {
	cfg := &RedisConfig{}
	viper.AddConfigPath(cl.ConfigPath)
	viper.SetConfigName(cl.ConfigName)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	isNotFoundError := func(m string) bool {
		return strings.Contains(strings.ToLower(m), "not found")
	}
	err := viper.ReadInConfig()
	if err != nil && !isNotFoundError(err.Error()) {
		return cfg, err
	}
	if err != nil && isNotFoundError(err.Error()) {
		cfg.RedisAddress = viper.GetString("REDIS_ADDRESS")
		cfg.RedisPassword = viper.GetString("REDIS_PASSWORD")
		cfg.RedisDB = viper.GetInt("REDIS_DB")

		return cfg, nil
	}
	err = viper.Unmarshal(&cfg)

	return cfg, err
}
