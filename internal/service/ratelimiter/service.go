package ratelimiter

import (
	"sync"
	"time"
)

type Limiter struct {
	rate      time.Duration
	capacity  int
	remaining int
	LastEvent time.Time
	lock      sync.Mutex
}

func NewLeakyBucket(rate time.Duration, capacity int) *Limiter {
	return &Limiter{
		rate:      rate,
		capacity:  capacity,
		remaining: capacity,
		LastEvent: time.Now(),
	}
}

func (lb *Limiter) Allow() bool {
	lb.lock.Lock()
	defer lb.lock.Unlock()

	now := time.Now()
	elapsed := now.Sub(lb.LastEvent)

	leakedTokens := int(elapsed / lb.rate)
	lb.remaining = min(lb.capacity, lb.remaining+leakedTokens)
	lb.LastEvent = now

	if lb.remaining > 0 {
		lb.remaining--
		return true
	}

	return false
}
