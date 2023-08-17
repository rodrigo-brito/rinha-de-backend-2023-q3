package main

import (
	"net/http"

	"rinha/etc/grpc"
	api "rinha/pkg/server/http"
)

func main() {
	storageClient := grpc.NewStorageProtobufClient("http://localhost:8080", http.DefaultClient)
	server := api.NewServer(storageClient)
	server.Start(8080)
}
