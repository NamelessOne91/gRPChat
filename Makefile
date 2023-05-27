PROJECT_DIR=${CURDIR}

# generate protobuf and grpc code
proto_gen:
	protoc chat.proto --proto_path=${PROJECT_DIR}/chat --go_out=. --go-grpc_out=.