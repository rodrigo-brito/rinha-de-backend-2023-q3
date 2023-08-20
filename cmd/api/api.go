package main

import (
	"fmt"
	"log"
	"net/http"

	"rinha/etc/grpc"
	"rinha/pkg/config"
	api "rinha/pkg/server/http"
)

func main() {
	storageClient := grpc.NewStorageProtobufClient(config.StorageAddress(), http.DefaultClient)
	server := api.NewServer(storageClient)
	fmt.Println("API running at http://localhot:8000")
	err := server.Start(8000, config.AccessLog())
	if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
