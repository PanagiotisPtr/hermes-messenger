package main

import (
	"context"
	"fmt"
	"net"

	"github.com/panagiotisptr/hermes-messenger/friends/server"
	"github.com/panagiotisptr/hermes-messenger/friends/server/connection"
	"github.com/panagiotisptr/hermes-messenger/friends/server/friends"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/grpcclientutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/grpcserviceutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/loggingutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/redisutils"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Provides the GRPC server instance
func ProvideGRPCServer(
	fs *server.FriendsServer,
	cfg *grpcserviceutils.GRPCServiceConfig,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterFriendsServer(gs, fs)

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
	logger *zap.Logger,
) {
	logger.Sugar().Info("Starting user service")
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
	app := fx.New(
		fx.Provide(
			utils.ProvideConfigLocation,
			grpcserviceutils.ProvideGRPCServiceConfig,
			loggingutils.ProvideProductionLogger,

			grpcclientutils.ProvideUserServiceClientConfig,
			grpcclientutils.ProvideUserServiceClient,

			mongoutils.ProvideMongoConfig,
			mongoutils.ProvideMongoClient,
			mongoutils.ProvideMongoDatabase,

			redisutils.ProvideRedisConfig,
			redisutils.ProvideRedisClient,

			server.ProvideFriendsServer,
			friends.ProvideFriendsService,
			connection.ProvideMongoRepository,
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
