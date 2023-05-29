PROJECT_DIR=${CURDIR}
CLIENT_BINARY=gRPChat-client
SERVER_BINARY=gRPChat-server

.PHONY:proto_gen build_client build_server

# generate protobuf and grpc code
proto_gen:
	@echo "Generating gRPC implementation"
	protoc chat.proto --proto_path=${PROJECT_DIR}/chat --go_out=. --go-grpc_out=.

build_client:
	@echo "Building client binary"
	cd ./client && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${CLIENT_BINARY} .

build_server:
	@echo "Building server binary"
	cd ./server && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${SERVER_BINARY} .