package main

import (
	"log"

	"github.com/panagiotisptr/hermes-messenger/libs/service"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/messaging/server"
	"github.com/panagiotisptr/hermes-messenger/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	healthCheckPort := utils.GetEnvVariableInt("HEALTH_CHECK_PORT", 12345)

	ipAddress, err := utils.GetMachineIpAddress()
	if err != nil {
		panic(err)
	}

	grpcService := service.NewGRPCService()
	err = grpcService.Bootstrap(service.GRPCServiceConfig{
		ServiceName:     "messaging",
		HostName:        ipAddress,
		ListenPort:      listenPort,
		HealthCheckPort: healthCheckPort,
	}, func(gs *grpc.Server, logger *log.Logger) error {
		cs := server.NewMessagingServer(logger)
		if err != nil {
			return err
		}
		protos.RegisterMessagingServer(gs, cs)
		reflection.Register(gs)

		return nil
	})
	if err != nil {
		panic(err)
	}

	err = grpcService.Start()
	if err != nil {
		panic(err)
	}
	defer grpcService.Stop()
}
