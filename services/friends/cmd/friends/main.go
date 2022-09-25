package main

import (
	"context"
	"fmt"
	"net"

	"github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/friends/config"
	"github.com/panagiotisptr/hermes-messenger/friends/server"
	"github.com/panagiotisptr/hermes-messenger/friends/server/connection"
	"github.com/panagiotisptr/hermes-messenger/friends/server/friends"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
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
	cfg *config.Config,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterFriendsServer(gs, fs)

	if cfg.GRPCReflection {
		reflection.Register(gs)
	}

	return gs, nil
}

func ProvideRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(cfg.Redis)
}

func ProvideMongoConfig(
	cfg *config.Config,
) *mongoutils.MongoConfig {
	return &cfg.MongoConfig
}

// Bootstraps the application
func Bootstrap(
	lc fx.Lifecycle,
	gs *grpc.Server,
	cfg *config.Config,
	logger *zap.Logger,
) {
	logger.Sugar().Info("Starting user service")
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

// Provides user Client
func ProvideUserClient(
	lc fx.Lifecycle,
	cfg *config.Config,
) (protos.UserClient, error) {
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

	return protos.NewUserClient(userConn), nil
}

func main() {
	app := fx.New(
		fx.Provide(
			ProvideRedisClient,
			ProvideLogger,
			ProvideUserClient,
			ProvideMongoConfig,
			mongoutils.ProvideMongoClient,
			mongoutils.ProvideMongoDatabase,
			config.ProvideConfig,
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
