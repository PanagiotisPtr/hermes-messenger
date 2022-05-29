package service

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
	"google.golang.org/grpc"
)

type GRPCServiceConfig struct {
	ServiceName     string
	HostName        string
	ListenPort      int
	HealthCheckPort int
}

type GRPCService struct {
	grpcServer     *grpc.Server
	client         *api.Client
	connectService *connect.Service
	listener       *net.Listener
	httpServer     *http.Server
}

func NewGRPCService() *GRPCService {
	return &GRPCService{
		grpcServer:     nil,
		client:         nil,
		connectService: nil,
		listener:       nil,
		httpServer:     nil,
	}
}

type SetupFunc func(*grpc.Server, *log.Logger) error

func (s *GRPCService) Bootstrap(config GRPCServiceConfig, setup SetupFunc) error {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return err
	}
	s.client = client

	connectService, err := connect.NewService(config.ServiceName, client)
	if err != nil {
		return err
	}
	s.connectService = connectService

	logger := log.New(os.Stdout, config.ServiceName+"-logs: ", log.Lshortfile)
	registration := &api.AgentServiceRegistration{
		ID:      config.ServiceName,
		Name:    config.ServiceName,
		Port:    config.ListenPort,
		Address: config.HostName,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%v/check", config.HostName, config.HealthCheckPort),
			Interval: "10s",
			Timeout:  "30s",
		},
	}

	err = s.client.Agent().ServiceRegister(registration)
	if err != nil {
		return err
	} else {
		logger.Printf("Successfully registered service: %s\n", config.ServiceName)
	}

	gs := grpc.NewServer()
	err = setup(gs, logger)
	if err != nil {
		return err
	}
	s.grpcServer = gs

	listenAddress := fmt.Sprintf(":%d", config.ListenPort)
	healthCheckAddress := fmt.Sprintf(":%d", config.HealthCheckPort)
	list, err := tls.Listen("tcp", listenAddress, s.connectService.ServerTLSConfig())
	if err != nil {
		return err
	}
	s.listener = &list

	mux := http.NewServeMux()

	mux.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Consul check")
	})

	s.httpServer = &http.Server{
		Addr:    healthCheckAddress,
		Handler: mux,
	}

	return nil
}

func (s *GRPCService) Start() error {
	err := s.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	err = s.grpcServer.Serve(*s.listener)

	return err
}

func (s *GRPCService) Stop() (error, error) {
	connectErr := s.connectService.Close()
	listError := (*s.listener).Close()

	return connectErr, listError
}
