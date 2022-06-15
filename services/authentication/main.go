package main

import (
	"context"
	"fmt"
	"net"

	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/authentication"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/secret"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/token"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/user"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// The application config
type Config struct {
	ListenPort     int
	GRPCReflection bool
}

// Provides the application config
func ProvideConfig() *Config {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)

	return &Config{
		ListenPort:     listenPort,
		GRPCReflection: grpcReflection,
	}
}

// Provides the GRPC server instance
func ProvideGRPCServer(
	as *server.AuthenticationServer,
	config *Config,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterAuthenticationServer(gs, as)

	if config.GRPCReflection {
		reflection.Register(gs)
	}

	return gs, nil
}

// Bootstraps the application
func Bootstrap(lc fx.Lifecycle, gs *grpc.Server, config *Config, logger *zap.Logger) {
	logger.Sugar().Info("Starting application")
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Info("Starting GRPC server.")
			addr := fmt.Sprintf(":%d", config.ListenPort)
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
			ProvideLogger,
			ProvideConfig,
			server.ProvideAuthenticationServer,
			authentication.ProvideAuthenticationService,
			token.ProvideTokenRepository,
			secret.ProvideMemorySecretRepository,
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
