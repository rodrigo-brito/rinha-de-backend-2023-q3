package main

import (
	"fmt"
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
	fmt.Println("Storage running at http://localhot:9000")
	err = server.Start(9000)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}
