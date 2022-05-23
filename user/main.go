package main

import (
	"hermes-messenger/user/protos"
	"hermes-messenger/user/server"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := log.New(os.Stdout, "messaging-service", log.Lshortfile)

	gs := grpc.NewServer()
	cs, err := server.NewUserService(logger)
	if err != nil {
		panic(err)
	}

	protos.RegisterUserServer(gs, cs)

	reflection.Register(gs)

	list, err := net.Listen("tcp", ":6060")
	if err != nil {
		logger.Fatalf("Unable to listen. Error: %v", err)
		os.Exit(1)
	} else {
		logger.Println("Listening on port :6060")
	}

	gs.Serve(list)
}
