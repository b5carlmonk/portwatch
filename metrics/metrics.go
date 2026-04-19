// Package metrics tracks runtime counters for portwatch scan cycles.
package metrics

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Counters holds cumulative scan metrics.
type Counters struct {
	mu          sync.Mutex
	Scans       int
	OpenPorts   int
	Changes     int
	Errors      int
	LastScan    time.Time
}

// Tracker records and reports scan metrics.
type Tracker struct {
	counters Counters
	w        io.Writer
}

// New returns a Tracker writing to stdout.
func New() *Tracker {
	return NewWithWriter(os.Stdout)
}

// NewWithWriter returns a Tracker writing to w.
func NewWithWriter(w io.Writer) *Tracker {
	return &Tracker{w: w}
}

// RecordScan increments scan count and records open port count and change count.
func (t *Tracker) RecordScan(openPorts, changes int, err error) {
	t.counters.mu.Lock()
	defer t.counters.mu.Unlock()
	t.counters.Scans++
	t.counters.OpenPorts = openPorts
	t.counters.Changes += changes
	t.counters.LastScan = time.Now()
	if err != nil {
		t.counters.Errors++
	}
}

// Snapshot returns a copy of the current counters.
func (t *Tracker) Snapshot() Counters {
	t.counters.mu.Lock()
	defer t.counters.mu.Unlock()
	return t.counters
}

// Print writes a summary of current metrics to the tracker's writer.
func (t *Tracker) Print() {
	s := t.Snapshot()
	fmt.Fprintf(t.w, "scans=%d open=%d changes=%d errors=%d last=%s\n",
		s.Scans, s.OpenPorts, s.Changes, s.Errors,
		s.LastScan.Format(time.RFC3339))
}
