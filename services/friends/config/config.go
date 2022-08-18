package config

import (
	"github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
)

// Application config
type Config struct {
	ListenPort         int
	GRPCReflection     bool
	UserServiceAddress string
	Redis              *redis.Options
}

// Provides the application config
func ProvideConfig() *Config {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	userServiceAddress := utils.GetEnvVariableString("USER_SERVICE_ADDR", "localhost:8080")
	redisAddress := utils.GetEnvVariableString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := utils.GetEnvVariableString("REDIS_PASSWORD", "")
	redisDatabase := utils.GetEnvVariableInt("REDIS_DB", 0)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)

	return &Config{
		ListenPort:         listenPort,
		GRPCReflection:     grpcReflection,
		UserServiceAddress: userServiceAddress,
		Redis: &redis.Options{
			Addr:     redisAddress,
			Password: redisPassword,
			DB:       redisDatabase,
		},
	}
}
