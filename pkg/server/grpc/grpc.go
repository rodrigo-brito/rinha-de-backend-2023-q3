package grpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/twitchtv/twirp"

	"rinha/etc/grpc"
	"rinha/pkg/entity"
	"rinha/pkg/service/search"
	"rinha/pkg/service/storage"
)

type Server struct {
	searcher search.Searcher
	storage  storage.Storage
}

func (s *Server) Search(ctx context.Context, request *grpc.SearchRequest) (*grpc.SearchResponse, error) {
	ids, err := s.searcher.Search(request.Term)
	if err != nil {
		return nil, err
	}

	response := grpc.SearchResponse{
		Users: make([]*grpc.User, 0, len(ids)),
	}

	for _, id := range ids {
		user, err := s.storage.Get(id)
		if err != nil {
			return nil, err
		}

		response.Users = append(response.Users, &grpc.User{
			Id:    id,
			Nick:  user.Nick,
			Name:  user.Name,
			Stack: user.Stack,
			Birth: user.Birth,
		})
	}

	return &response, nil
}

func (s *Server) Get(ctx context.Context, request *grpc.GetRequest) (*grpc.User, error) {
	user, err := s.storage.Get(request.Id)
	if err == storage.ErrNotFound {
		return &grpc.User{}, twirp.NewError(twirp.NotFound, err.Error())
	} else if err != nil {
		return &grpc.User{}, err
	}

	if user != nil {
		return &grpc.User{
			Id:    request.Id,
			Nick:  user.Nick,
			Name:  user.Name,
			Stack: user.Stack,
			Birth: user.Birth,
		}, nil
	}

	return nil, nil
}

func (s *Server) Total(context.Context, *grpc.TotalRequest) (*grpc.TotalResponse, error) {
	total := s.storage.Total()
	return &grpc.TotalResponse{
		Total: int64(total),
	}, nil
}

func (s *Server) Save(ctx context.Context, user *grpc.User) (*grpc.SaveResponse, error) {
	data := entity.User{
		Nick:  user.Nick,
		Name:  user.Name,
		Birth: user.Birth,
		Stack: user.Stack,
	}

	var err error
	data.ID, err = s.storage.Store(data)
	if err == storage.ErrDuplicated {
		return nil, twirp.NewError(twirp.AlreadyExists, err.Error())
	}
	go s.searcher.Save(data)

	return &grpc.SaveResponse{
		Id: data.ID,
	}, nil
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
