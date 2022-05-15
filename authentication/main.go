package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"panagiotisptr/authentication/protos"
	"panagiotisptr/authentication/server"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logger := log.New(os.Stdout, "messaging-service", log.Lshortfile)
	port := 80

	portEnvVariable, foundPortEnvVariable := os.LookupEnv("PORT")
	if !foundPortEnvVariable {
		fmt.Println("Missing PORT env variable")
		os.Exit(1)
	}
	port, err := strconv.Atoi(portEnvVariable)
	if err != nil {
		panic(err)
	}

	gs := grpc.NewServer()
	cs, err := server.NewAuthenticationServer(logger)
	if err != nil {
		panic(err)
	}

	protos.RegisterAuthenticationServer(gs, cs)

	reflection.Register(gs)

	address := fmt.Sprintf(":%d", port)
	list, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("Unable to listen. Error: %v", err)
		os.Exit(1)
	} else {
		logger.Println("Listening on port " + address)
	}

	gs.Serve(list)
}
