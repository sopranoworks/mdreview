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

func (s *Store) StartCleanup(interval time.Duration, maxAge time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			s.mu.Lock()
			now := time.Now()
			for id, p := range s.previews {
				if now.Sub(p.CreatedAt) > maxAge {
					delete(s.previews, id)
				}
			}
			s.mu.Unlock()
		}
	}()
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
