// Package shadow detects ports that are open but not present in a known
// baseline, flagging them as "shadow" services that may be unexpected.
package shadow

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/scanner"
)

// Entry describes a single shadow port finding.
type Entry struct {
	Host     string
	Port     int
	Protocol string
	Reason   string
}

// Detector compares live scan results against a known-good set and reports
// any ports that are not present in the allowed set.
type Detector struct {
	mu      sync.RWMutex
	allowed map[string]struct{}
}

// New returns a Detector pre-loaded with the provided baseline results.
func New(baseline []scanner.Result) *Detector {
	d := &Detector{
		allowed: make(map[string]struct{}, len(baseline)),
	}
	for _, r := range baseline {
		if r.Open {
			d.allowed[key(r)] = struct{}{}
		}
	}
	return d
}

// Detect returns all results from current that are open but absent from the
// baseline that was provided when the Detector was created.
func (d *Detector) Detect(current []scanner.Result) []Entry {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var findings []Entry
	for _, r := range current {
		if !r.Open {
			continue
		}
		if _, ok := d.allowed[key(r)]; !ok {
			findings = append(findings, Entry{
				Host:     r.Host,
				Port:     r.Port,
				Protocol: r.Protocol,
				Reason:   "port not present in baseline",
			})
		}
	}
	return findings
}

// Allow adds a result to the allowed set at runtime.
func (d *Detector) Allow(r scanner.Result) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.allowed[key(r)] = struct{}{}
}

// Len returns the number of entries in the allowed baseline.
func (d *Detector) Len() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.allowed)
}

func key(r scanner.Result) string {
	return fmt.Sprintf("%s:%s:%d", r.Host, r.Protocol, r.Port)
}
