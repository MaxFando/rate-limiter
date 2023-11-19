package inmemory

import (
	"context"
	"sync"

	"github.com/MaxFando/rate-limiter/internal/service/ratelimiter"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]*ratelimiter.Limiter
}

func New() *Store {
	return &Store{
		data: make(map[string]*ratelimiter.Limiter),
		mu:   sync.RWMutex{},
	}
}

func (db *Store) GetAll(ctx context.Context) map[string]*ratelimiter.Limiter {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return db.data
}

func (db *Store) Get(ctx context.Context, key string) (*ratelimiter.Limiter, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	item, ok := db.data[key]

	if !ok {
		return nil, false
	}

	return item, true
}

func (db *Store) Set(ctx context.Context, key string, value *ratelimiter.Limiter) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[key] = value
}

func (db *Store) Delete(ctx context.Context, key string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, key)
}
