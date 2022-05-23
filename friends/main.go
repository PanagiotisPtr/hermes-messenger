package main

import (
	"hermes-messenger/friends/protos"
	"hermes-messenger/friends/server"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := log.New(os.Stdout, "friends-service", log.Lshortfile)

	gs := grpc.NewServer()
	cs, err := server.NewFriendsServer(logger)
	if err != nil {
		panic(err)
	}

	protos.RegisterFriendsServer(gs, cs)

	reflection.Register(gs)

	list, err := net.Listen("tcp", ":7070")
	if err != nil {
		logger.Fatalf("Unable to listen. Error: %v", err)
		os.Exit(1)
	} else {
		logger.Println("Listening on port :7070")
	}

	gs.Serve(list)
}
