package config

import (
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
)

// Application config
type Config struct {
	ListenPort     int
	GRPCReflection bool
	Redis          *redis.Options
	ESConfig       elasticsearch.Config
	MongoConfig    mongoutils.MongoConfig
}

// Provides the application config
func ProvideConfig() *Config {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	redisAddress := utils.GetEnvVariableString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := utils.GetEnvVariableString("REDIS_PASSWORD", "")
	redisDatabase := utils.GetEnvVariableInt("REDIS_DB", 0)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)
	esAddresses := utils.GetEnvVariableString("ES_ADDRESSES", "http://localhost:9200")
	mongoUri := utils.GetEnvVariableString("MONGO_URI", "mongodb://localhost:27017")
	mongoDb := utils.GetEnvVariableString("MONGO_DB", "user")

	return &Config{
		ListenPort:     listenPort,
		GRPCReflection: grpcReflection,
		Redis: &redis.Options{
			Addr:     redisAddress,
			Password: redisPassword,
			DB:       redisDatabase,
		},
		ESConfig: elasticsearch.Config{
			Addresses: strings.Split(esAddresses, ","),
		},
		MongoConfig: mongoutils.MongoConfig{
			MongoUri: mongoUri,
			MongoDb:  mongoDb,
		},
	}
}
