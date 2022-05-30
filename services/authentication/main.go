package main

import (
	"log"

	"github.com/panagiotisptr/hermes-messenger/libs/service"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server"

	"google.golang.org/grpc"
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
		ServiceName:     "authentication",
		HostName:        ipAddress,
		ListenPort:      listenPort,
		HealthCheckPort: healthCheckPort,
		GRPCReflection:  true,
	}, func(gs *grpc.Server, logger *log.Logger) error {
		cs, err := server.NewAuthenticationServer(logger)
		if err != nil {
			return err
		}
		protos.RegisterAuthenticationServer(gs, cs)

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
