package grpcserviceutils

import (
	"strings"

	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/spf13/viper"
)

// ServiceConfig is the configuration required
// for the service runtime itself
type GRPCServiceConfig struct {
	ServicePort    int  `mapstructure:"SERVICE_PORT"`
	GRPCReflection bool `mapstructure:"GRPC_REFLECTION"`
}

func ProvideGRPCServiceConfig(cl *utils.ConfigLocation) (*GRPCServiceConfig, error) {
	cfg := &GRPCServiceConfig{}
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
		return cfg, nil
	}
	err = viper.Unmarshal(&cfg)

	return cfg, err
}
