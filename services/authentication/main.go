package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/panagiotisptr/hermes-messenger/libs/service"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func getIpAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("Could not find IP address of machine")
}

func main() {
	listenPort := 80
	listenPortEnvVariable, foundListenPortEnvVariable := os.LookupEnv("LISTEN_PORT")
	if foundListenPortEnvVariable {
		lp, err := strconv.Atoi(listenPortEnvVariable)
		if err != nil {
			panic(err)
		}
		listenPort = lp
	}

	healthCheckPort := 12345
	healthCheckPortEnvVariable, foundHealthCheckPortEnvVariable := os.LookupEnv("HEALTH_CHECK_PORT")
	if foundHealthCheckPortEnvVariable {
		hcp, err := strconv.Atoi(healthCheckPortEnvVariable)
		if err != nil {
			panic(err)
		}
		healthCheckPort = hcp
	}

	ipAddress, err := getIpAddress()
	if err != nil {
		panic(err)
	}

	grpcService := service.NewGRPCService()
	err = grpcService.Bootstrap(service.GRPCServiceConfig{
		ServiceName:     "authentication",
		HostName:        ipAddress,
		ListenPort:      listenPort,
		HealthCheckPort: healthCheckPort,
	}, func(gs *grpc.Server, logger *log.Logger) error {
		cs, err := server.NewAuthenticationServer(logger)
		if err != nil {
			return err
		}
		protos.RegisterAuthenticationServer(gs, cs)
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
