package main

import (
	"io"

	"github.com/NamelessOne91/grpc-chat/chat"
)

// Connection represents the gRPC client's connection to a chat server with a gRPC stream
type Connection struct {
	conn     chat.Chat_ChatServer
	sendChan chan *chat.ChatMessage
	quitChan chan struct{}
}

// NewConnection inits and returns a pointer to a new Connection for the given gRPC stream
func NewConnection(conn chat.Chat_ChatServer) *Connection {
	c := &Connection{
		conn:     conn,
		sendChan: make(chan *chat.ChatMessage),
		quitChan: make(chan struct{}),
	}
	go c.start()
	return c
}

// start is run in its own goroutine and pulls incoming broadcasted messages to be sent on the stream
func (c *Connection) start() {
	running := true
	for running {
		select {
		case msg := <-c.sendChan:
			// errors are ignored and result in the message not being received
			c.conn.Send(msg)
		case <-c.quitChan:
			running = false
		}
	}
}

// Close send a termination signal to the Connection, the returned error is always nil
func (c *Connection) Close() error {
	close(c.quitChan)
	close(c.sendChan)
	return nil
}

// Send sends a message on the Connection's dedicated channel
func (c *Connection) Send(msg *chat.ChatMessage) {
	defer func() {
		// ignore panics for sending on closed channels
		recover()
	}()
	c.sendChan <- msg
}

// GetMessages pulls incoming messages, through the gRPC stream, and sends them on the passed channel for broadcasting.
//
// Runs untill one of the following conditions is true:
//   - the client closes the stream due to the user quitting
//   - the gRPC chat server is no longer reachable
//   - another error occurs
func (c *Connection) GetMessages(broadcastChan chan<- *chat.ChatMessage) error {
	for {
		msg, err := c.conn.Recv()
		if err == io.EOF {
			c.Close()
			return nil
		} else if err != nil {
			c.Close()
			return err
		}

		go func(msg *chat.ChatMessage) {
			select {
			case broadcastChan <- msg:
			case <-c.quitChan:
			}
		}(msg)
	}
}
