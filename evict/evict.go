// Package evict provides a time-based eviction tracker that removes stale
// port scan results which have not been seen within a configurable TTL window.
package evict

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/scanner"
)

// entry holds the last time a result was observed.
type entry struct {
	result  scanner.Result
	lastSeen time.Time
}

// Tracker records the last-seen timestamp for each port result and can
// evict entries that have exceeded the TTL.
type Tracker struct {
	mu      sync.Mutex
	entries map[string]entry
	ttl     time.Duration
	now     func() time.Time
}

// New creates a Tracker with the given TTL duration.
func New(ttl time.Duration) *Tracker {
	return &Tracker{
		entries: make(map[string]entry),
		ttl:     ttl,
		now:     time.Now,
	}
}

// Observe records or refreshes the last-seen time for each result.
func (t *Tracker) Observe(results []scanner.Result) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	for _, r := range results {
		t.entries[key(r)] = entry{result: r, lastSeen: now}
	}
}

// Evict removes entries older than the TTL and returns the evicted results.
func (t *Tracker) Evict() []scanner.Result {
	t.mu.Lock()
	defer t.mu.Unlock()
	cutoff := t.now().Add(-t.ttl)
	var evicted []scanner.Result
	for k, e := range t.entries {
		if e.lastSeen.Before(cutoff) {
			evicted = append(evicted, e.result)
			delete(t.entries, k)
		}
	}
	return evicted
}

// Active returns all results currently tracked (not yet evicted).
func (t *Tracker) Active() []scanner.Result {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]scanner.Result, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, e.result)
	}
	return out
}

// Len returns the number of tracked entries.
func (t *Tracker) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}

func key(r scanner.Result) string {
	return fmt.Sprintf("%s:%d/%s", r.Host, r.Port, r.Proto)
}
