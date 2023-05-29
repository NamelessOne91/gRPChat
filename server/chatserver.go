package main

import (
	"sync"

	"github.com/NamelessOne91/grpc-chat/chat"
)

// ChatServer represents a gRPC chat server that can hold multiple connections to clients
// and broadcasts messages to all the connected clients
type ChatServer struct {
	chat.UnimplementedChatServer
	broadcastChan chan *chat.ChatMessage
	quitChan      chan struct{}
	connections   []*Connection
	m             sync.Mutex
}

// NewChatServer returns a pointer to an initialized and started ChatServer
func NewChatServer() *ChatServer {
	srv := &ChatServer{
		broadcastChan: make(chan *chat.ChatMessage),
		quitChan:      make(chan struct{}),
	}
	go srv.start()
	return srv
}

// start is run in its own goroutine, it instructs the server to wait for incoming messages
// and concurrently broadcast them to all connected clients
func (c *ChatServer) start() {
	running := true
	for running {
		select {
		case msg := <-c.broadcastChan:
			// lock the connection
			c.m.Lock()
			// concurrently send the message to all connected clients
			for _, v := range c.connections {
				go v.Send(msg)
			}
			// release connection
			c.m.Unlock()
		case <-c.quitChan:
			running = false
		}
	}
}

// Close sends a termination signal to the ChatServer, the returned error is always nil
func (c *ChatServer) Close() error {
	close(c.quitChan)
	return nil
}

// Chat is the gRPC Server API implementation called in a separate goroutine for each new stream opened by a client.
//
// It will keep a pointer to the underlying gRPC connection and keep pulling new messages to be broadcasted.
//
// Once the connection has been closed, it will be removed from the slice of the target connections for broadcasting.
func (c *ChatServer) Chat(stream chat.Chat_ChatServer) error {
	conn := NewConnection(stream)

	// register new connection
	c.m.Lock()
	c.connections = append(c.connections, conn)
	c.m.Unlock()

	// blocks untill the connection is closed
	err := conn.GetMessages(c.broadcastChan)

	c.m.Lock()
	for i, v := range c.connections {
		if v == conn {
			c.connections = append(c.connections[:i], c.connections[i+1:]...)
		}
	}
	c.m.Unlock()

	return err
}
