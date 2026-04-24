// Package cooldown enforces a minimum quiet period between repeated alerts
// for the same key. Once an alert fires, further alerts for that key are
// suppressed until the cooldown duration has elapsed.
package cooldown

import (
	"sync"
	"time"
)

// Cooldown tracks the last-fired time for each key.
type Cooldown struct {
	mu       sync.Mutex
	duration time.Duration
	lastFired map[string]time.Time
	now       func() time.Time
}

// New returns a Cooldown that enforces the given minimum duration between
// successive events for the same key.
func New(d time.Duration) *Cooldown {
	return &Cooldown{
		duration:  d,
		lastFired: make(map[string]time.Time),
		now:       time.Now,
	}
}

// Allow returns true if the key has not fired within the cooldown window.
// When true is returned the last-fired timestamp for the key is updated.
func (c *Cooldown) Allow(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	if last, ok := c.lastFired[key]; ok {
		if now.Sub(last) < c.duration {
			return false
		}
	}
	c.lastFired[key] = now
	return true
}

// Reset removes the cooldown record for a single key, allowing the next
// call to Allow for that key to succeed immediately.
func (c *Cooldown) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.lastFired, key)
}

// Flush removes all cooldown records.
func (c *Cooldown) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastFired = make(map[string]time.Time)
}

// Remaining returns how long until the key's cooldown expires.
// Returns 0 if the key is not in cooldown.
func (c *Cooldown) Remaining(key string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	last, ok := c.lastFired[key]
	if !ok {
		return 0
	}
	remaining := c.duration - c.now().Sub(last)
	if remaining < 0 {
		return 0
	}
	return remaining
}
