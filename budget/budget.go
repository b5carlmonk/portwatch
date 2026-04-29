// Package budget limits the total number of scan targets processed
// within a rolling time window, protecting the host from excessive load.
package budget

import (
	"errors"
	"sync"
	"time"
)

// ErrBudgetExceeded is returned when the scan budget has been exhausted.
var ErrBudgetExceeded = errors.New("scan budget exceeded for current window")

// Entry records a single budget consumption event.
type Entry struct {
	At    time.Time
	Count int
}

// Budget tracks cumulative scan target usage within a sliding window.
type Budget struct {
	mu      sync.Mutex
	max     int
	window  time.Duration
	entries []Entry
	now     func() time.Time
}

// New creates a Budget that allows at most max targets per window duration.
func New(max int, window time.Duration) *Budget {
	return &Budget{
		max:    max,
		window: window,
		now:    time.Now,
	}
}

// newWithClock creates a Budget with an injectable clock for testing.
func newWithClock(max int, window time.Duration, clock func() time.Time) *Budget {
	return &Budget{max: max, window: window, now: clock}
}

// Allow attempts to consume n units from the budget.
// It returns ErrBudgetExceeded if adding n would surpass the limit.
func (b *Budget) Allow(n int) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.evict()

	total := 0
	for _, e := range b.entries {
		total += e.Count
	}

	if total+n > b.max {
		return ErrBudgetExceeded
	}

	b.entries = append(b.entries, Entry{At: b.now(), Count: n})
	return nil
}

// Remaining returns the number of units still available in the current window.
func (b *Budget) Remaining() int {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.evict()

	used := 0
	for _, e := range b.entries {
		used += e.Count
	}

	r := b.max - used
	if r < 0 {
		return 0
	}
	return r
}

// Reset clears all recorded usage.
func (b *Budget) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.entries = nil
}

// evict removes entries that have fallen outside the rolling window.
// Must be called with b.mu held.
func (b *Budget) evict() {
	cutoff := b.now().Add(-b.window)
	i := 0
	for i < len(b.entries) && b.entries[i].At.Before(cutoff) {
		i++
	}
	b.entries = b.entries[i:]
}
