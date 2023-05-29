package main

import (
	"flag"
	"log"
	"net"

	"github.com/NamelessOne91/grpc-chat/chat"
	"google.golang.org/grpc"
)

func main() {
	// parse command line args
	webPort := flag.String("p", "8080", "TCP port to use")
	flag.Parse()

	lst, err := net.Listen("tcp", ":"+*webPort)
	if err != nil {
		log.Fatalf("cannot listen on port %s: %v", *webPort, err)
	}

	s := grpc.NewServer()
	srv := NewChatServer()
	chat.RegisterChatServer(s, srv)

	log.Printf("gRPChat server started and listening on port %s\n", *webPort)
	err = s.Serve(lst)
	if err != nil {
		log.Fatalf("received error, shutting down server: %v", err)
	}
}
