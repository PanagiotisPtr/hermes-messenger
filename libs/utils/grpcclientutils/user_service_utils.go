package grpcclientutils

import (
	"context"
	"strings"

	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

type UserServiceClientConfig struct {
	UserServiceAddress string `mapstructure:"USER_SERVICE_ADDRESS"`
}

func ProvideUserServiceClientConfig(cl *utils.ConfigLocation) (*UserServiceClientConfig, error) {
	cfg := &UserServiceClientConfig{}
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
		cfg.UserServiceAddress = viper.GetString("USER_SERVICE_ADDRESS")

		return cfg, nil
	}
	err = viper.Unmarshal(&cfg)

	return cfg, err
}

// Provides user Client
func ProvideUserServiceClient(
	lc fx.Lifecycle,
	cfg *UserServiceClientConfig,
) (protos.UserServiceClient, error) {
	userConn, err := grpc.Dial(
		cfg.UserServiceAddress,
		grpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return userConn.Close()
		},
	})

	return protos.NewUserServiceClient(userConn), nil
}
