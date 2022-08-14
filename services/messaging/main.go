package main

import (
	"context"
	"fmt"
	"net"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/panagiotisptr/hermes-messenger/messaging/config"
	"github.com/panagiotisptr/hermes-messenger/messaging/server"
	"github.com/panagiotisptr/hermes-messenger/messaging/server/messaging"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Provides the GRPC server instance
func ProvideGRPCServer(
	ms *server.MessagingServer,
	cfg *config.Config,
) (*grpc.Server, error) {
	gs := grpc.NewServer()
	protos.RegisterMessagingServer(gs, ms)

	if cfg.GRPCReflection {
		reflection.Register(gs)
	}

	return gs, nil
}

func ProvideElasticsearchClient(cfg *config.Config) (*elasticsearch.Client, error) {
	return elasticsearch.NewClient(cfg.ESConfig)
}

// Bootstraps the application
func Bootstrap(lc fx.Lifecycle, gs *grpc.Server, cfg *config.Config, logger *zap.Logger) {
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
			ProvideElasticsearchClient,
			ProvideLogger,
			config.ProvideConfig,
			server.ProvideUserServer,
			messaging.ProvideMessagingService,
			messaging.ProvideESRepository,
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
