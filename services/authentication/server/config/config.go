package config

import (
	"time"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
)

// The application config
type Config struct {
	UUID                      uuid.UUID
	ListenPort                int
	GRPCReflection            bool
	RedisAddress              string
	RedisPassword             string
	RedisDatabase             int
	RefreshTokenDuration      time.Duration
	AccessTokenDuration       time.Duration
	KeyPairGenerationInterval time.Duration
}

// Provides the application config
func ProvideConfig() *Config {
	secondsInDay := 60 * 60 * 24
	secondsInHour := 60 * 60
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	redisAddress := utils.GetEnvVariableString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := utils.GetEnvVariableString("REDIS_PASSWORD", "")
	redisDatabase := utils.GetEnvVariableInt("REDIS_DB", 0)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)
	refreshTokenDurationSeconds := utils.GetEnvVariableInt("REFRESH_TOKEN_EXP_SEC", secondsInDay)
	accessTokenDurationSeconds := utils.GetEnvVariableInt("ACCESS_TOKEN_EXP_SEC", secondsInHour)
	keyPairGenerationInterval := utils.GetEnvVariableInt("KEY_PAIR_GENERATION_INTERVAL", secondsInDay)

	return &Config{
		UUID:                      uuid.New(),
		ListenPort:                listenPort,
		GRPCReflection:            grpcReflection,
		RedisAddress:              redisAddress,
		RedisPassword:             redisPassword,
		RedisDatabase:             redisDatabase,
		RefreshTokenDuration:      time.Second * time.Duration(refreshTokenDurationSeconds),
		AccessTokenDuration:       time.Second * time.Duration(accessTokenDurationSeconds),
		KeyPairGenerationInterval: time.Second * time.Duration(keyPairGenerationInterval),
	}
}
