// Package throttle limits how frequently scan cycles can be triggered
// to prevent resource exhaustion from rapid successive calls.
package throttle

import (
	"sync"
	"time"
)

// Throttle enforces a minimum interval between allowed actions.
type Throttle struct {
	mu       sync.Mutex
	interval time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New returns a Throttle that enforces the given minimum interval per key.
func New(interval time.Duration) *Throttle {
	return &Throttle{
		interval: interval,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if enough time has passed since the last allowed call
// for the given key. If allowed, it records the current time.
func (t *Throttle) Allow(key string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	if last, ok := t.last[key]; ok {
		if now.Sub(last) < t.interval {
			return false
		}
	}
	t.last[key] = now
	return true
}

// Reset clears the recorded time for a specific key.
func (t *Throttle) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, key)
}

// Flush clears all recorded times.
func (t *Throttle) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = make(map[string]time.Time)
}

// Remaining returns how much time is left before the key is allowed again.
// Returns zero if the key is already allowed.
func (t *Throttle) Remaining(key string) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	last, ok := t.last[key]
	if !ok {
		return 0
	}
	elapsed := t.now().Sub(last)
	if elapsed >= t.interval {
		return 0
	}
	return t.interval - elapsed
}
