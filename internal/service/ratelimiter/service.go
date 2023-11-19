package ratelimiter

import (
	"golang.org/x/time/rate"
	"time"
)

type Limiter struct {
	rate      *rate.Limiter
	LastEvent time.Time
}

func NewLimiter(r rate.Limit, b int) *Limiter {
	limiter := rate.NewLimiter(r, b)
	return &Limiter{rate: limiter}
}

func (t *Limiter) Allow() bool {
	t.LastEvent = time.Now()
	allow := t.rate.Allow()
	return allow
}
