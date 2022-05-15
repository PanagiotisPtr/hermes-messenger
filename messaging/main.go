package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"panagiotisptr/messaging/protos"
	"panagiotisptr/messaging/server"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	uuid1 := uuid.New()
	uuid2 := uuid.New()

	uuid3bytes := make([]byte, 16)
	for index, _ := range uuid1 {
		uuid3bytes[index] = uuid1[index] ^ uuid2[index]
	}
	uuid3, _ := uuid.FromBytes(uuid3bytes)

	uuid4bytes := make([]byte, 16)
	for index, _ := range uuid2 {
		uuid4bytes[index] = uuid2[index] ^ uuid1[index]
	}
	uuid4, _ := uuid.FromBytes(uuid3bytes)

	fmt.Println(uuid1.String())
	fmt.Println(uuid2.String())
	fmt.Println(uuid3.String())
	fmt.Println(uuid4.String())

	logger := log.New(os.Stdout, "messaging-service", log.Lshortfile)

	gs := grpc.NewServer()
	cs := server.NewMessagingServer(logger)

	protos.RegisterMessagingServer(gs, cs)

	reflection.Register(gs)

	list, err := net.Listen("tcp", ":9090")
	if err != nil {
		logger.Fatalf("Unable to listen. Error: %v", err)
		os.Exit(1)
	} else {
		logger.Println("Listening on port :9090")
	}

	gs.Serve(list)
}
