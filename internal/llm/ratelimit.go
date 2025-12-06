package llm

import (
	"sync"
	"time"
)

// RateLimiter tracks request rates per user.
type RateLimiter interface {
	// Allow checks if a request is allowed for the given user.
	// Returns true if allowed, false if rate limited.
	Allow(userID string) bool
}

// InMemoryRateLimiter implements a sliding window rate limiter.
type InMemoryRateLimiter struct {
	maxRequests int
	window      time.Duration

	mu       sync.Mutex
	requests map[string][]time.Time
}

// NewInMemoryRateLimiter creates a rate limiter that allows maxRequests per window.
func NewInMemoryRateLimiter(maxRequests int, window time.Duration) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make(map[string][]time.Time),
	}
}

// Allow checks if a request is allowed for the given user.
func (r *InMemoryRateLimiter) Allow(userID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	// Get existing requests and filter out old ones
	existing := r.requests[userID]
	var recent []time.Time
	for _, t := range existing {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}

	// Check if under limit
	if len(recent) >= r.maxRequests {
		r.requests[userID] = recent
		return false
	}

	// Record this request
	recent = append(recent, now)
	r.requests[userID] = recent
	return true
}

// Cleanup removes stale entries from the rate limiter.
// Call this periodically to prevent memory growth.
func (r *InMemoryRateLimiter) Cleanup() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-r.window)

	for userID, times := range r.requests {
		var recent []time.Time
		for _, t := range times {
			if t.After(windowStart) {
				recent = append(recent, t)
			}
		}
		if len(recent) == 0 {
			delete(r.requests, userID)
		} else {
			r.requests[userID] = recent
		}
	}
}
