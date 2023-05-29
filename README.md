# gRPChat
A simple gRPC based CLI chat, using bidirectional streams - written in Go

Includes both server and client implementations.
The server supports multiple connections and broadcasts messages to all connected users.

## Setup

To build the executable binaries, **on Linux**, you can use the provided Makefile commands:
 - position yourself in the project's root folder
 - run `make build_client` to build the client's executable binary
 - run `make build_server` to build the server's executable binary

 Binaries will be put in the **/client/bin/** and **/server/bin/** folders.

## Running the server

The server's executable can be run without providing any flag.
This results in the server listening on port *8080*

You can set the following flags when launching the server:

| Flag | Type | Default | Meaning
| :---:|:--:|:--:|:--|
| p | string | 8080 | TCP port on which the server will listen |

## Running the client

In order to run the client you **must** provide a username with the `-name` flag.
By default, the client will try to connect to a gRPC server on localhost:8080. 
You can change the server's URL using the `-url` flag.

| Flag | Type | Default | Meaning
| :---:|:--:|:--:|:--|
| url | string | localhost:8080 | URL of the gRPC server |
| name | string | / | username to be displayed

Once the client is running, press the **Enter** key to send what's written in your terminal as message to the server.

To quit, send the message `quit` or press `CTRL + C`