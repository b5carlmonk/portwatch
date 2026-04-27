// Package window provides a sliding time-window counter for tracking
// scan events and computing rates over a configurable duration.
package window

import (
	"sync"
	"time"
)

// Entry holds a single timestamped count.
type Entry struct {
	At    time.Time
	Count int
}

// Window accumulates entries and exposes totals within a sliding window.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	entries  []Entry
	now      func() time.Time
}

// New returns a Window that retains entries within the given duration.
func New(d time.Duration) *Window {
	return &Window{
		duration: d,
		now:      time.Now,
	}
}

// newWithClock is used in tests to inject a fixed clock.
func newWithClock(d time.Duration, fn func() time.Time) *Window {
	return &Window{duration: d, now: fn}
}

// Add records n events at the current time.
func (w *Window) Add(n int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = append(w.entries, Entry{At: w.now(), Count: n})
	w.evict()
}

// Total returns the sum of all counts within the window.
func (w *Window) Total() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	total := 0
	for _, e := range w.entries {
		total += e.Count
	}
	return total
}

// Len returns the number of entries currently in the window.
func (w *Window) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	return len(w.entries)
}

// Reset removes all entries.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = nil
}

// evict removes entries older than the window duration. Caller must hold mu.
func (w *Window) evict() {
	cutoff := w.now().Add(-w.duration)
	i := 0
	for i < len(w.entries) && w.entries[i].At.Before(cutoff) {
		i++
	}
	w.entries = w.entries[i:]
}
