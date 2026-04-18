// Package ratelimit provides a simple token-bucket rate limiter for
// suppressing repeated alerts about the same port event.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-key event times and suppresses events that occur
// more frequently than the configured interval.
type Limiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     map[string]time.Time
}

// New returns a Limiter that allows at most one event per key per interval.
func New(interval time.Duration) *Limiter {
	return &Limiter{
		interval: interval,
		last:     make(map[string]time.Time),
	}
}

// Allow returns true if the key has not been seen within the interval.
// If allowed, the timestamp for the key is updated.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.interval {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the recorded timestamp for a key, allowing the next
// event through immediately.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// Flush removes all recorded timestamps.
func (l *Limiter) Flush() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
