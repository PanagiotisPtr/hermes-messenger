package config

import (
	"io/ioutil"
	"strings"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
)

// Application config
type Config struct {
	ListenPort     int
	GRPCReflection bool
	ESConfig       elasticsearch.Config
}

// Provides the application config
func ProvideConfig() (*Config, error) {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	esAddresses := utils.GetEnvVariableString("ES_ADDRESSES", "https://localhost:8200")
	esUsername := utils.GetEnvVariableString("ES_USERNAME", "elastic")
	esPassword := utils.GetEnvVariableString("ES_PASSWORD", "")
	esCertPath := utils.GetEnvVariableString("ES_CERT_PATH", "config/certs/http_ca.crt")
	esCertFingerprint := utils.GetEnvVariableString("ES_CERT_FINGERPRINT", "")
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)

	esCert, err := ioutil.ReadFile(esCertPath)
	if err != nil {
		return nil, err
	}

	return &Config{
		ListenPort:     listenPort,
		GRPCReflection: grpcReflection,
		ESConfig: elasticsearch.Config{
			Addresses:              strings.Split(esAddresses, ","),
			Username:               esUsername,
			Password:               esPassword,
			CACert:                 esCert,
			CertificateFingerprint: esCertFingerprint,
		},
	}, nil
}
