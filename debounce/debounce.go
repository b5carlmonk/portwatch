// Package debounce delays alert delivery until a change persists for a
// configurable number of consecutive scan cycles, reducing noise from
// transient port flaps.
package debounce

import (
	"sync"

	"github.com/user/portwatch/scanner"
)

// Entry tracks how many consecutive cycles a given change has been seen.
type Entry struct {
	Count    int
	LastDiff scanner.Diff
}

// Debouncer holds pending changes and only confirms them once they have been
// observed for at least Threshold consecutive cycles.
type Debouncer struct {
	mu        sync.Mutex
	threshold int
	pending   map[string]*Entry
}

// New returns a Debouncer that requires threshold consecutive observations
// before a change is considered confirmed.
func New(threshold int) *Debouncer {
	if threshold < 1 {
		threshold = 1
	}
	return &Debouncer{
		threshold: threshold,
		pending:   make(map[string]*Entry),
	}
}

// Evaluate accepts a Diff and returns only the changes that have been
// observed for at least Threshold consecutive cycles. Changes that have not
// yet reached the threshold are held internally. Keys that disappear from
// the diff are reset.
func (d *Debouncer) Evaluate(diff scanner.Diff) scanner.Diff {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Build a set of keys present in this diff.
	current := make(map[string]bool)

	confirmed := scanner.Diff{
		Opened: []scanner.Result{},
		Closed: []scanner.Result{},
	}

	process := func(results []scanner.Result, tag string) []scanner.Result {
		var out []scanner.Result
		for _, r := range results {
			k := tag + ":" + r.Host + ":" + itoa(r.Port) + ":" + r.Proto
			current[k] = true
			e, ok := d.pending[k]
			if !ok {
				e = &Entry{}
				d.pending[k] = e
			}
			e.Count++
			e.LastDiff = diff
			if e.Count >= d.threshold {
				out = append(out, r)
				delete(d.pending, k)
			}
		}
		return out
	}

	confirmed.Opened = process(diff.Opened, "open")
	confirmed.Closed = process(diff.Closed, "close")

	// Evict stale pending keys that were not present in this diff.
	for k := range d.pending {
		if !current[k] {
			delete(d.pending, k)
		}
	}

	return confirmed
}

// PendingLen returns the number of changes currently awaiting confirmation.
func (d *Debouncer) PendingLen() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.pending)
}

// Flush discards all pending (unconfirmed) change entries.
func (d *Debouncer) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.pending = make(map[string]*Entry)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
