package main

import (
	"hermes-messenger/config/protos"
	"hermes-messenger/config/server"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := log.New(os.Stdout, "messaging-service", log.Lshortfile)

	gs := grpc.NewServer()
	cs, err := server.NewConfigServer(logger)
	if err != nil {
		panic(err)
	}

	protos.RegisterConfigServer(gs, cs)

	reflection.Register(gs)

	list, err := net.Listen("tcp", ":8888")
	if err != nil {
		logger.Fatalf("Unable to listen. Error: %v", err)
		os.Exit(1)
	} else {
		logger.Println("Listening on port :8888")
	}

	gs.Serve(list)
}
