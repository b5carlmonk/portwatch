// Package decay implements a time-based score decay tracker.
// Scores assigned to port/host keys decrease toward zero over time,
// allowing the system to "forget" stale risk signals automatically.
package decay

import (
	"math"
	"sync"
	"time"
)

// Entry holds a score and the last time it was updated.
type Entry struct {
	Score     float64
	UpdatedAt time.Time
}

// Tracker applies exponential decay to named keys.
type Tracker struct {
	mu       sync.Mutex
	entries  map[string]Entry
	halfLife time.Duration
	clock    func() time.Time
}

// New creates a Tracker with the given half-life duration.
// After one half-life, a score is halved.
func New(halfLife time.Duration) *Tracker {
	return newWithClock(halfLife, time.Now)
}

func newWithClock(halfLife time.Duration, clock func() time.Time) *Tracker {
	return &Tracker{
		entries:  make(map[string]Entry),
		halfLife: halfLife,
		clock:    clock,
	}
}

// Add increases the score for key by delta after applying decay.
func (t *Tracker) Add(key string, delta float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	e, ok := t.entries[key]
	if !ok {
		t.entries[key] = Entry{Score: delta, UpdatedAt: now}
		return
	}
	decayed := t.decayScore(e.Score, e.UpdatedAt, now)
	t.entries[key] = Entry{Score: decayed + delta, UpdatedAt: now}
}

// Get returns the current decayed score for key, or 0 if not found.
func (t *Tracker) Get(key string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	if !ok {
		return 0
	}
	return t.decayScore(e.Score, e.UpdatedAt, t.clock())
}

// Reset removes the entry for key.
func (t *Tracker) Reset(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, key)
}

// Flush removes all tracked entries.
func (t *Tracker) Flush() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[string]Entry)
}

// decayScore computes score * 0.5^(elapsed / halfLife).
func (t *Tracker) decayScore(score float64, since, now time.Time) float64 {
	elapsed := now.Sub(since)
	if elapsed <= 0 {
		return score
	}
	exponent := float64(elapsed) / float64(t.halfLife)
	return score * math.Pow(0.5, exponent)
}
