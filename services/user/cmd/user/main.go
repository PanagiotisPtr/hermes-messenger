package main

import (
	"context"
	"fmt"
	"net"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis/v9"
	"github.com/panagiotisptr/hermes-messenger/libs/utils/mongoutils"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/user/config"
	"github.com/panagiotisptr/hermes-messenger/user/server"
	"github.com/panagiotisptr/hermes-messenger/user/server/user"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Provides the GRPC server instance
func ProvideGRPCServer(
	us *server.UserServer,
	cfg *config.Config,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterUserServer(gs, us)

	if cfg.GRPCReflection {
		reflection.Register(gs)
	}

	return gs, nil
}

func ProvideElasticsearchClient(cfg *config.Config) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(cfg.ESConfig)
}

func ProvideRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(cfg.Redis)
}

func ProvideMongoClient(
	lc fx.Lifecycle,
	logger *zap.Logger,
	cfg *config.Config,
) (*mongo.Client, error) {
	client, err := mongo.NewClient(
		mongoutils.SetRegistryForUuids(
			options.Client().ApplyURI(cfg.MongoConfig.MongoUri),
		),
	)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Sugar().Info("Connecting to mongo")
			err := client.Connect(ctx)
			if err != nil {
				return err
			}

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Sugar().Info("Disconnecting from database")

			return client.Disconnect(ctx)
		},
	})

	return client, nil
}

func ProvideMongoDatabase(
	client *mongo.Client,
	cfg *config.Config,
) *mongo.Database {
	return client.Database(cfg.MongoConfig.MongoDb)
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
			ProvideLogger,
			ProvideMongoClient,
			ProvideMongoDatabase,
			config.ProvideConfig,
			server.ProvideUserServer,
			user.ProvideUserService,
			user.ProvideMongoRepository,
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