package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"

	"github.com/NamelessOne91/grpc-chat/chat"
	"google.golang.org/grpc"
)

// ChatClient represents a single client connected to the gRPC chat server.
//
// Holds info about the chosen username and a pointer to the underlying bidirectional gRPC stream
type ChatClient struct {
	stream   chat.Chat_ChatClient
	waitChan chan struct{}
	username string
}

// NewChatClient inits the bidirectional gRPC stream and returns a pointer to a new ChatClient.
//
// Messages sent by the returned client will appear as being sent by the provided username
func NewChatClient(conn *grpc.ClientConn, username string) *ChatClient {
	client := chat.NewChatClient(conn)
	stream, err := client.Chat(context.Background())
	if err != nil {
		log.Fatalf("failed to connect to init chat client: %v", err)
	}

	cc := &ChatClient{
		stream:   stream,
		waitChan: make(chan struct{}),
		username: username,
	}
	go cc.start()

	return cc
}

// start is run in its own goroutine, it instructs the client to wait for user inputs
// and send them on the gRPC stream or close the connection and exit the program
func (c *ChatClient) start() {
	log.Println("Connected - type \"quit\" or press CTRL + C to exit")
	// read CLI input
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		// disconnect client
		if msg == "quit" {
			err := c.stream.CloseSend()
			c.Close()
			if err != nil {
				log.Fatalf("error while disconnecting from the chat server: %v", err)
			}
			break
		}

		// send along the gRPC stream to server
		err := c.stream.Send(&chat.ChatMessage{
			User:    c.username,
			Message: msg,
		})
		if err != nil {
			log.Fatalf("error while sending message: %v", err)
		}
	}
}

// Close sends a termination signal to the ChatClient, the returned error is always nil
func (c *ChatClient) Close() error {
	close(c.waitChan)
	return nil
}

// Chat instructs the client to pull incoming messages, through the gRPC stream, from the server and display them.
//
// Runs untill one of the following conditions is true:
//   - the client closes the stream due to the user quitting
//   - the gRPC chat server is no longer reachable
//   - another error occurs
func (c *ChatClient) Chat() error {
	for {
		msg, err := c.stream.Recv()
		if errors.Is(err, io.EOF) {
			// disconnected client
			break
		} else if err != nil {
			// TODO: error handling of other cases
			c.Close()
			return errors.New("the connection with gRPChat server has been lost")
		}
		log.Printf("%s: %s\n", msg.User, msg.Message)
	}
	<-c.waitChan
	return nil
}
