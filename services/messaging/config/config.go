package config

import (
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	redis "github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
)

// Application config
type Config struct {
	ListenPort            int
	GRPCReflection        bool
	ESConfig              elasticsearch.Config
	FriendsServiceAddress string
	Redis                 *redis.Options
	MongoConfig           mongoutils.MongoConfig
}

// Provides the application config
func ProvideConfig() *Config {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	friendsServiceAddres := utils.GetEnvVariableString("FRIENDS_SERVICE_ADDR", "localhost:8080")
	esAddresses := utils.GetEnvVariableString("ES_ADDRESSES", "https://localhost:9200")
	esUsername := utils.GetEnvVariableString("ES_USERNAME", "elastic")
	esPassword := utils.GetEnvVariableString("ES_PASSWORD", "")
	redisAddress := utils.GetEnvVariableString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := utils.GetEnvVariableString("REDIS_PASSWORD", "")
	redisDatabase := utils.GetEnvVariableInt("REDIS_DB", 0)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)
	mongoUri := utils.GetEnvVariableString("MONGO_URI", "mongodb://localhost:27017")
	mongoDb := utils.GetEnvVariableString("MONGO_DB", "messaging")

	return &Config{
		ListenPort:            listenPort,
		GRPCReflection:        grpcReflection,
		FriendsServiceAddress: friendsServiceAddres,
		ESConfig: elasticsearch.Config{
			Addresses: strings.Split(esAddresses, ","),
			Username:  esUsername,
			Password:  esPassword,
		},
		Redis: &redis.Options{
			Addr:     redisAddress,
			Password: redisPassword,
			DB:       redisDatabase,
		},
		MongoConfig: mongoutils.MongoConfig{
			MongoUri: mongoUri,
			MongoDb:  mongoDb,
		},
	}
}
