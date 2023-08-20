package storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"sync"

	"rinha/pkg/entity"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrDuplicated = errors.New("duplicated")
)

type Storage interface {
	Get(id string) (*entity.User, error)
	Store(entity.User) (string, error)
	Total() int
}

type storage struct {
	sync.Mutex
	data map[string]entity.User
}

func (s *storage) Get(id string) (*entity.User, error) {
	user, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &user, nil
}

func (s *storage) UUID(nick string) string {
	key := fmt.Sprintf("%x", md5.Sum([]byte(nick)))
	return fmt.Sprintf("%s-%s-%s-%s-%s", key[:8], key[8:12], key[12:16], key[16:20], key[20:32])
}

func (s *storage) Store(user entity.User) (string, error) {
	user.ID = s.UUID(user.Nick)
	if _, ok := s.data[user.ID]; ok {
		return user.ID, ErrDuplicated
	}
	s.Lock()
	defer s.Unlock()
	s.data[user.ID] = user
	return user.ID, nil
}

func (s *storage) Total() int {
	return len(s.data)
}

func NewStorage() *storage {
	return &storage{
		data: make(map[string]entity.User),
	}
}
