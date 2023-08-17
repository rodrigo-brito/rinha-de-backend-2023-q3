package grpc

import (
	"context"
	"fmt"
	"net/http"

	"rinha/etc/grpc"
	"rinha/pkg/service/search"
	"rinha/pkg/service/storage"
)

type Server struct {
	searcher search.Searcher
	storage  storage.Storage
}

func (s *Server) Search(ctx context.Context, request *grpc.SearchRequest) (*grpc.Response, error) {
	return nil, nil
}

func (s *Server) Get(ctx context.Context, request *grpc.GetRequest) (*grpc.Response, error) {
	return nil, nil
}

func NewServer(searcher search.Searcher, storage storage.Storage) *Server {
	return &Server{
		searcher: searcher,
		storage:  storage,
	}
}

func (s *Server) Start(port int) error {
	handler := grpc.NewStorageServer(s)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
