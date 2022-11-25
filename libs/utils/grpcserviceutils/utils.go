package grpcserviceutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/spf13/viper"
	"google.golang.org/grpc/metadata"
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
		cfg.ServicePort = viper.GetInt("SERVICE_PORT")
		cfg.GRPCReflection = viper.GetBool("GRPC_REFLECTION")

		return cfg, nil
	}
	err = viper.Unmarshal(&cfg)

	return cfg, err
}

func LoadMetadataValuesToContext(ctx context.Context, keys ...string) (context.Context, error) {
	m, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, fmt.Errorf("missing grpc metadata in request")
	}
	for _, k := range keys {
		values := m.Get(k)
		if len(values) == 0 {
			return ctx, fmt.Errorf("missing key '%s' in request metadata", k)
		}
		ctx = context.WithValue(ctx, k, values[0])
	}

	return ctx, nil
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userId, ok := ctx.Value("user-id").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("missing userId in context")
	}

	return userId, nil
}

func WithUserID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(
		ctx,
		"user-id",
		id,
	)
}
