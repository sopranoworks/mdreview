package server

import (
	"sync"
	"time"
)

type Preview struct {
	Content   string
	CreatedAt time.Time
}

type Store struct {
	mu       sync.RWMutex
	previews map[string]Preview
}

func NewStore() *Store {
	return &Store{
		previews: make(map[string]Preview),
	}
}

func (s *Store) Set(id string, content string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.previews[id] = Preview{
		Content:   content,
		CreatedAt: time.Now(),
	}
}

func (s *Store) Get(id string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.previews[id]
	return p.Content, ok
}
