package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/grpcserviceutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/loggingutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/rabbitmqutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/redisutils"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server/authentication"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Provides the GRPC server instance
func ProvideGRPCServer(
	as *server.AuthenticationServer,
	cfg *grpcserviceutils.GRPCServiceConfig,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterAuthenticationServer(gs, as)

	if cfg.GRPCReflection {
		reflection.Register(gs)
	}

	return gs, nil
}

// Bootstraps the application
func Bootstrap(
	lc fx.Lifecycle,
	gs *grpc.Server,
	cfg *grpcserviceutils.GRPCServiceConfig,
	service *authentication.Service,
	logger *zap.Logger,
) {
	logger.Sugar().Info("Starting application")
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Info("Starting GRPC server.")
			addr := fmt.Sprintf(":%d", cfg.ServicePort)
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

func main() {
	flag.Parse()
	app := fx.New(
		fx.Provide(
			loggingutils.ProvideProductionLogger,
			utils.ProvideConfigLocation,
			grpcserviceutils.ProvideGRPCServiceConfig,
			redisutils.ProvideRedisConfig,
			redisutils.ProvideRedisClient,
			rabbitmqutils.ProvideAmqpConfig,
			rabbitmqutils.ProvideAmqpChannel,
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
