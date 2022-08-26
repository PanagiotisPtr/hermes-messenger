package config

import (
	"io/ioutil"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
)

// Application config
type Config struct {
	ListenPort     int
	GRPCReflection bool
	Redis          *redis.Options
	ESConfig       elasticsearch.Config
}

// Provides the application config
func ProvideConfig() (*Config, error) {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	redisAddress := utils.GetEnvVariableString("REDIS_ADDRESS", "localhost:6379")
	redisPassword := utils.GetEnvVariableString("REDIS_PASSWORD", "")
	redisDatabase := utils.GetEnvVariableInt("REDIS_DB", 0)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)
	esAddresses := utils.GetEnvVariableString("ES_ADDRESSES", "https://localhost:8200")
	esUsername := utils.GetEnvVariableString("ES_USERNAME", "elastic")
	esPassword := utils.GetEnvVariableString("ES_PASSWORD", "")
	esCertPath := utils.GetEnvVariableString("ES_CERT_PATH", "config/certs/http_ca.crt")
	esCertFingerprint := utils.GetEnvVariableString("ES_CERT_FINGERPRINT", "")

	esCert, err := ioutil.ReadFile(esCertPath)
	if err != nil {
		return nil, err
	}

	return &Config{
		ListenPort:     listenPort,
		GRPCReflection: grpcReflection,
		Redis: &redis.Options{
			Addr:     redisAddress,
			Password: redisPassword,
			DB:       redisDatabase,
		},
		ESConfig: elasticsearch.Config{
			Addresses:              strings.Split(esAddresses, ","),
			Username:               esUsername,
			Password:               esPassword,
			CACert:                 esCert,
			CertificateFingerprint: esCertFingerprint,
		},
	}, nil
}
