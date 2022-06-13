package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/panagiotisptr/hermes-messenger/libs/utils"
	"github.com/panagiotisptr/hermes-messenger/protos"
	"github.com/panagiotisptr/hermes-messenger/services/authentication/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listenPort := utils.GetEnvVariableInt("LISTEN_PORT", 80)
	grpcReflection := utils.GetEnvVariableBool("GRPC_REFLECTION", false)

	logger := log.New(os.Stdout, "[authentication-logs] ", log.Lshortfile)
	gs := grpc.NewServer()
	cs, err := server.NewAuthenticationServer(logger)
	if err != nil {
		panic(err)
	}
	protos.RegisterAuthenticationServer(gs, cs)
	if err != nil {
		panic(err)
	}

	if grpcReflection {
		reflection.Register(gs)
	}

	addr := fmt.Sprintf(":%d", listenPort)
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	} else {
		logger.Println("Listening on " + addr)
	}

	gs.Serve(list)
}
