package main

import (
	"log"

	"github.com/panagiotisptr/hermes-messenger/friends/server"
	"github.com/panagiotisptr/hermes-messenger/libs/service"
	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/protos"

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
		ServiceName:     "friends",
		HostName:        ipAddress,
		ListenPort:      listenPort,
		HealthCheckPort: healthCheckPort,
		GRPCReflection:  true,
	}, func(gs *grpc.Server, logger *log.Logger) error {
		cs, err := server.NewFriendsServer(logger)
		if err != nil {
			return err
		}
		protos.RegisterFriendsServer(gs, cs)

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
