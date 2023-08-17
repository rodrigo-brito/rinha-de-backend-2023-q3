package main

import (
	"log"

	"rinha/pkg/server/grpc"
	"rinha/pkg/service/search"
	"rinha/pkg/service/storage"
)

func main() {
	storageService := storage.NewStorage()
	searchService, err := search.NewSearcher()
	if err != nil {
		log.Fatalf("error creating searcher: %s", err)
	}

	server := grpc.NewServer(searchService, storageService)
	err = server.Start(8080)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
