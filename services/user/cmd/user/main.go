package main

import (
	"context"
	"fmt"
	"net"

	"github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/user/config"
	mongo_repository "github.com/panagiotisptr/hermes-messenger/user/repository/mongo"
	"github.com/panagiotisptr/hermes-messenger/user/server"
	"github.com/panagiotisptr/hermes-messenger/user/service"

	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Provides the GRPC server instance
func ProvideGRPCServer(
	us *server.Server,
	cfg *config.Config,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterUserServiceServer(gs, us)

	if cfg.Service.GRPCReflection {
		reflection.Register(gs)
	}

	return gs, nil
}

// ProvideRedisClient Provides Redis Client
func ProvideRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(cfg.Redis)
}

// ProvideMongoConfig provides mongo config
func ProvideMongoConfig(
	cfg *config.Config,
) *mongoutils.Config {
	return &cfg.Mongo
}

// Bootstraps the application
func Bootstrap(
	lc fx.Lifecycle,
	client *mongo.Client,
	gs *grpc.Server,
	cfg *config.Config,
	logger *zap.Logger,
) {
	logger.Sugar().Info("Starting user service")
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Info("Starting GRPC server.")
			addr := fmt.Sprintf(":%d", cfg.Service.Port)
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
			ProvideMongoConfig,
			mongoutils.ProvideMongoClient,
			mongoutils.ProvideMongoDatabase,
			config.ProvideConfig,
			server.ProvideUserServer,
			service.ProvideUserService,
			mongo_repository.ProvideUserRepository,
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
