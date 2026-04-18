// Package digest generates periodic summary reports of port scan activity.
package digest

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/scanner"
)

// Entry holds a timestamped scan result snapshot.
type Entry struct {
	Time    time.Time
	Results []scanner.Result
}

// Digest accumulates entries and writes summaries.
type Digest struct {
	w       io.Writer
	entries []Entry
}

// New returns a Digest writing to stdout.
func New() *Digest {
	return &Digest{w: os.Stdout}
}

// NewWithWriter returns a Digest writing to w.
func NewWithWriter(w io.Writer) *Digest {
	return &Digest{w: w}
}

// Add appends a new entry to the digest.
func (d *Digest) Add(results []scanner.Result) {
	d.entries = append(d.entries, Entry{
		Time:    time.Now(),
		Results: results,
	})
}

// Flush writes a summary of all accumulated entries and clears them.
func (d *Digest) Flush() {
	if len(d.entries) == 0 {
		fmt.Fprintln(d.w, "[digest] no data collected")
		return
	}
	fmt.Fprintf(d.w, "[digest] summary over %d snapshots\n", len(d.entries))
	portSeen := map[string]int{}
	for _, e := range d.entries {
		for _, r := range e.Results {
			for _, p := range r.Ports {
				key := fmt.Sprintf("%d/%s", p.Port, p.Proto)
				portSeen[key]++
			}
		}
	}
	for k, count := range portSeen {
		fmt.Fprintf(d.w, "  port %-12s seen in %d snapshot(s)\n", k, count)
	}
	d.entries = nil
}

// Len returns the number of accumulated entries.
func (d *Digest) Len() int { return len(d.entries) }
