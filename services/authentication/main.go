package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/google/uuid"
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

type Config struct {
	ListenPort     int
	GRPCReflection bool
}

func NewLogger() (*log.Logger, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	logger := log.New(os.Stdout, "[authentication]["+uuid.String()+"]", log.Lshortfile)

	return logger, nil
}

func ProvideConfig() *Config {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)

	return &Config{
		ListenPort:     listenPort,
		GRPCReflection: grpcReflection,
	}
}

func ProvideGRPCServer(
	as *server.AuthenticationServer,
	config *Config,
	logger *log.Logger,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterAuthenticationServer(gs, as)

	if config.GRPCReflection {
		reflection.Register(gs)
	}

	return gs, nil
}

func Bootstrap(lc fx.Lifecycle, gs *grpc.Server, config *Config, slogger *zap.SugaredLogger) {
	slogger.Info("Starting application")
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			slogger.Info("Starting GRPC server.")
			addr := fmt.Sprintf(":%d", config.ListenPort)
			list, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			} else {
				slogger.Info("Listening on " + addr)
			}

			go gs.Serve(list)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			slogger.Info("Stopping GRPC server.")
			gs.Stop()

			return slogger.Sync()
		},
	})
}

func ProvideLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	slogger := logger.Sugar()

	return slogger
}

func main() {
	app := fx.New(
		fx.Provide(
			NewLogger,
			ProvideLogger,
			ProvideConfig,
			server.NewAuthenticationServer,
			authentication.NewService,
			token.NewRepository,
			secret.NewMemoryRepository,
			user.NewMemoryRepository,
			ProvideGRPCServer,
		),
		fx.Invoke(Bootstrap),
		fx.WithLogger(
			func() fxevent.Logger {
				return fxevent.NopLogger
			},
		),
	)

	app.Run()
}
