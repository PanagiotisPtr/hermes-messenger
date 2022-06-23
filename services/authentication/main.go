package main

import (
	"context"
	"fmt"
	"net"

	redis "github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/authentication"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/config"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/keys"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/token"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/user"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Provides the GRPC server instance
func ProvideGRPCServer(
	as *server.AuthenticationServer,
	cfg *config.Config,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterAuthenticationServer(gs, as)

	if cfg.GRPCReflection {
		reflection.Register(gs)
	}

	return gs, nil
}

func ProvideRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDatabase,
	})
}

// Bootstraps the application
func Bootstrap(lc fx.Lifecycle, gs *grpc.Server, cfg *config.Config, logger *zap.Logger) {
	logger.Sugar().Info("Starting application")
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Info("Starting GRPC server.")
			addr := fmt.Sprintf(":%d", cfg.ListenPort)
			list, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			} else {
				logger.Sugar().Info("Listening on " + addr)
			}

			go gs.Serve(list)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Sugar().Info("Stopping GRPC server.")
			gs.Stop()

			return logger.Sync()
		},
	})
}

// Provides the ZAP logger
func ProvideLogger() *zap.Logger {
	logger, _ := zap.NewProduction()

	return logger
}

func main() {
	app := fx.New(
		fx.Provide(
			ProvideRedisClient,
			ProvideLogger,
			config.ProvideConfig,
			server.ProvideAuthenticationServer,
			authentication.ProvideAuthenticationService,
			token.ProvideTokenRepository,
			keys.ProvideRedisKeysRepository,
			user.ProvideMemoryUserRepository,
			ProvideGRPCServer,
		),
		fx.Invoke(Bootstrap),
		fx.WithLogger(
			func(logger *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: logger}
			},
		),
	)

	app.Run()
}
