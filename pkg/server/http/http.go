package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"rinha/etc/grpc"
)

type Server struct {
	storageClient grpc.Storage
}

func NewServer(storageClient grpc.Storage) *Server {
	return &Server{
		storageClient: storageClient,
	}
}

func (s *Server) Start(port int) {
	r := chi.NewRouter()
	r.Post("/pessoas", s.CreateHandle)
	r.Post("/pessoas/:id", s.GetHandle)
	r.Post("/pessoas", s.SearchHandle)
	r.Post("/contagem-pessoas", s.CountHandle)
}

func (s *Server) GetHandle(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) SearchHandle(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) CountHandle(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) CreateHandle(w http.ResponseWriter, r *http.Request) {
}
