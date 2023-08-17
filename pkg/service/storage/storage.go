package storage

import (
	"errors"

	"rinha/pkg/entity"
)

var ErrNotFound = errors.New("not found")

type Storage interface {
	Get(id string) (*entity.People, error)
	Store(entity.People)
}

type storage struct {
	data map[string]entity.People
}

func (s storage) Get(id string) (*entity.People, error) {
	people, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &people, nil
}

func (s storage) Store(people entity.People) {
	s.data[people.Nick] = people
}

func NewStorage() *storage {
	return &storage{
		data: make(map[string]entity.People),
	}
}
