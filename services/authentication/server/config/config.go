package config

import (
	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
)

// The application config
type Config struct {
	UUID           uuid.UUID
	ListenPort     int
	GRPCReflection bool
	RedisAddress   string
	RedisPassword  string
	RedisDatabase  int
}

// Provides the application config
func ProvideConfig() *Config {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	redisAddress := utils.GetEnvVariableString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := utils.GetEnvVariableString("REDIS_PASSWORD", "")
	redisDatabase := utils.GetEnvVariableInt("REDIS_DB", 0)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)

	return &Config{
		UUID:           uuid.New(),
		ListenPort:     listenPort,
		GRPCReflection: grpcReflection,
		RedisAddress:   redisAddress,
		RedisPassword:  redisPassword,
		RedisDatabase:  redisDatabase,
	}
}
