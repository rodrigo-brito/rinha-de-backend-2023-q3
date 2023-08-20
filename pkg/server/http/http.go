package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/twitchtv/twirp"

	"rinha/etc/grpc"
	"rinha/pkg/entity"
)

type Server struct {
	storageClient grpc.Storage
}

func NewServer(storageClient grpc.Storage) *Server {
	return &Server{
		storageClient: storageClient,
	}
}

func (s *Server) Start(port int, accessLog bool) error {
	r := chi.NewRouter()
	if accessLog {
		r.Use(middleware.DefaultLogger)
	}
	r.Post("/pessoas", s.CreateHandle)
	r.Get("/pessoas/{id}", s.GetHandle)
	r.Get("/pessoas", s.SearchHandle)
	r.Get("/contagem-pessoas", s.CountHandle)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func (s *Server) GetHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := s.storageClient.Get(r.Context(), &grpc.GetRequest{Id: id})
	if err != nil {
		if twerr, ok := err.(twirp.Error); ok {
			if twerr.Code() == twirp.NotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}

		slog.Error(fmt.Sprintf("error getting user: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if user != nil {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(entity.User{
			ID:    user.Id,
			Nick:  user.Nick,
			Name:  user.Name,
			Birth: user.Birth,
			Stack: user.Stack,
		})
		if err != nil {
			slog.Error(fmt.Sprintf("error encoding user: %s", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func (s *Server) SearchHandle(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("t")
	if term == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := s.storageClient.Search(r.Context(), &grpc.SearchRequest{Term: term})
	if err != nil {
		slog.Error(fmt.Sprintf("error getting user: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result := make([]entity.User, 0, len(users.Users))
	for _, user := range users.Users {
		result = append(result, entity.User{
			ID:    user.Id,
			Nick:  user.Nick,
			Name:  user.Name,
			Birth: user.Birth,
			Stack: user.Stack,
		})
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		slog.Error(fmt.Sprintf("error encoding user: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) CountHandle(w http.ResponseWriter, r *http.Request) {
	result, err := s.storageClient.Total(r.Context(), &grpc.TotalRequest{})
	if err != nil {
		slog.Error(fmt.Sprintf("error getting user: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("%d", result.Total)))
	if err != nil {
		slog.Error(fmt.Sprintf("error encoding user: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) CreateHandle(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		slog.Error(fmt.Sprintf("error decoding user: %s", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !user.Validate() {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	id, err := s.storageClient.Save(r.Context(), &grpc.User{
		Nick:  user.Nick,
		Name:  user.Name,
		Stack: user.Stack,
		Birth: user.Birth,
	})
	if err != nil {
		if twerr, ok := err.(twirp.Error); ok {
			if twerr.Code() == twirp.AlreadyExists {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		}

		slog.Error(fmt.Sprintf("error creating user: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		slog.Error(fmt.Sprintf("error encoding user: %s", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
