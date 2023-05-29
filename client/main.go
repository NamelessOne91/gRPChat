package main

import (
	"flag"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// parse CLI flags
	serverUrl := flag.String("url", "localhost:8080", "URL of the gRPC chat server - default: localhost:8080")
	username := flag.String("name", "", "the username associated to your messages")
	flag.Parse()

	if *username == "" {
		log.Fatal("please, provide an username with the -name flag")
	}

	// connect to the gRPC chat server
	conn, err := grpc.Dial(*serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC chat server: %v", err)
	}
	defer conn.Close()

	// init client: receive and send messages
	client := NewChatClient(conn, *username)
	err = client.Chat()
	if err != nil {
		log.Fatalf("client stopped with error: %v", err)
	}
}
