// Package quota enforces per-host scan rate limits based on a maximum
// number of allowed scans within a rolling time window.
package quota

import (
	"fmt"
	"sync"
	"time"
)

// Entry tracks scan timestamps for a single key.
type Entry struct {
	times []time.Time
}

// Quota enforces a maximum number of events within a sliding window.
type Quota struct {
	mu      sync.Mutex
	entries map[string]*Entry
	max     int
	window  time.Duration
}

// New creates a Quota that allows at most max events per window per key.
func New(max int, window time.Duration) *Quota {
	return &Quota{
		entries: make(map[string]*Entry),
		max:     max,
		window:  window,
	}
}

// Allow returns true if the key is within quota, recording the attempt.
// It returns false and an error if the quota has been exceeded.
func (q *Quota) Allow(key string) (bool, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-q.window)

	e, ok := q.entries[key]
	if !ok {
		e = &Entry{}
		q.entries[key] = e
	}

	// Evict timestamps outside the window.
	valid := e.times[:0]
	for _, t := range e.times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	e.times = valid

	if len(e.times) >= q.max {
		return false, fmt.Errorf("quota exceeded for %q: %d/%d scans in %s", key, len(e.times), q.max, q.window)
	}

	e.times = append(e.times, now)
	return true, nil
}

// Reset clears the quota state for a specific key.
func (q *Quota) Reset(key string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.entries, key)
}

// Flush clears all quota state.
func (q *Quota) Flush() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.entries = make(map[string]*Entry)
}

// Remaining returns how many more events are allowed for key within the window.
func (q *Quota) Remaining(key string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	cutoff := time.Now().Add(-q.window)
	e, ok := q.entries[key]
	if !ok {
		return q.max
	}
	count := 0
	for _, t := range e.times {
		if t.After(cutoff) {
			count++
		}
	}
	rem := q.max - count
	if rem < 0 {
		return 0
	}
	return rem
}
