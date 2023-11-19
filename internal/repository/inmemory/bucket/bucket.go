package bucket

import (
	"context"
	"time"

	"github.com/MaxFando/rate-limiter/internal/service/ratelimiter"
	"github.com/MaxFando/rate-limiter/internal/store/inmemory"
)

type Repository struct {
	storage *inmemory.Store
}

func NewRepository(storage *inmemory.Store) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) DeleteUnusedBucket(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	for {
		<-ticker.C
		for ip, limiter := range r.storage.GetAll(ctx) {
			if time.Since(limiter.LastEvent) > time.Duration(60)*time.Second {
				r.storage.Delete(ctx, ip)
			}
		}
	}
}

func (r *Repository) TryGetPermissionInBucket(ctx context.Context, requestValue string, limit int) bool {
	limiter, ok := r.storage.Get(ctx, requestValue)
	if !ok {
		newLimiter := r.newBucket(limit)
		r.storage.Set(ctx, requestValue, newLimiter)

		return newLimiter.Allow()
	}

	return limiter.Allow()
}

func (r *Repository) ResetBucket(ctx context.Context, requestValue string) bool {
	_, ok := r.storage.Get(ctx, requestValue)
	if !ok {
		return false
	}
	r.storage.Delete(ctx, requestValue)
	return true
}

func (r *Repository) newBucket(limit int) *ratelimiter.Limiter {
	limiter := ratelimiter.NewLeakyBucket(time.Duration(60)*time.Second, limit)
	return limiter
}
