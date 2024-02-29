package middleware

import (
	"sync"
	"time"
)

type RateLimiter struct {
	rate     float64
	capacity float64
	tokens   float64
	mutex    sync.Mutex
}

func NewRateLimiter(rate float64) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		capacity: rate,
		tokens:   rate,
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := float64(time.Now().UnixNano()) / float64(time.Second)

	elapsed := now - rl.tokens
	if elapsed > 0 {
		rl.tokens += elapsed * rl.rate
		if rl.tokens > rl.capacity {
			rl.tokens = rl.capacity
		}
	}
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}

	return false
}
