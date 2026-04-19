// Package rollup aggregates multiple scan diffs into a single summary
// over a configurable time window.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/scanner"
)

// Entry holds a diff captured at a point in time.
type Entry struct {
	At   time.Time
	Diff scanner.Diff
}

// Rollup collects diffs and can summarise them.
type Rollup struct {
	mu      sync.Mutex
	window  time.Duration
	entries []Entry
}

// New returns a Rollup that keeps entries within the given window.
func New(window time.Duration) *Rollup {
	return &Rollup{window: window}
}

// Add appends a diff, pruning entries older than the window.
func (r *Rollup) Add(d scanner.Diff) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.entries = append(r.entries, Entry{At: now, Diff: d})
	r.prune(now)
}

func (r *Rollup) prune(now time.Time) {
	cutoff := now.Add(-r.window)
	var keep []Entry
	for _, e := range r.entries {
		if e.At.After(cutoff) {
			keep = append(keep, e)
		}
	}
	r.entries = keep
}

// Summary returns aggregated opened/closed port counts within the window.
func (r *Rollup) Summary() (opened, closed int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prune(time.Now())
	for _, e := range r.entries {
		opened += len(e.Diff.Opened)
		closed += len(e.Diff.Closed)
	}
	return
}

// Entries returns a snapshot of current entries.
func (r *Rollup) Entries() []Entry {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prune(time.Now())
	out := make([]Entry, len(r.entries))
	copy(out, r.entries)
	return out
}

// Flush clears all entries.
func (r *Rollup) Flush() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = nil
}
