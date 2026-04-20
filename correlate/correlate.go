// Package correlate groups related port change events across multiple hosts
// into correlated incidents based on timing and protocol similarity.
package correlate

import (
	"sync"
	"time"

	"github.com/user/portwatch/scanner"
)

// Event represents a single port change detected on a host.
type Event struct {
	Host      string
	Port      int
	Proto     string
	State     string // "opened" or "closed"
	Timestamp time.Time
}

// Incident groups related events that occurred within a correlation window.
type Incident struct {
	ID     string
	Events []Event
}

// Correlator collects events and groups them into incidents.
type Correlator struct {
	mu      sync.Mutex
	window  time.Duration
	events  []Event
	closed  func(Incident)
}

// New returns a Correlator that groups events within the given time window.
// The onIncident callback is invoked when an incident is flushed.
func New(window time.Duration, onIncident func(Incident)) *Correlator {
	return &Correlator{
		window:     window,
		onIncident: onIncident,
	}
}

// Add records a new event derived from a scanner diff entry.
func (c *Correlator) Add(host string, port int, proto, state string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.events = append(c.events, Event{
		Host:      host,
		Port:      port,
		Proto:     proto,
		State:     state,
		Timestamp: time.Now(),
	})
}

// Flush groups all buffered events that fall within the correlation window
// into incidents and fires the callback for each, then clears the buffer.
func (c *Correlator) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.events) == 0 {
		return
	}
	incidents := group(c.events, c.window)
	for _, inc := range incidents {
		if c.onIncident != nil {
			c.onIncident(inc)
		}
	}
	c.events = nil
}

// group partitions events into incidents by proximity in time.
func group(events []Event, window time.Duration) []Incident {
	if len(events) == 0 {
		return nil
	}
	var incidents []Incident
	current := Incident{ID: incidentID(events[0]), Events: []Event{events[0]}}
	for _, e := range events[1:] {
		ref := current.Events[len(current.Events)-1].Timestamp
		if e.Timestamp.Sub(ref) <= window {
			current.Events = append(current.Events, e)
		} else {
			incidents = append(incidents, current)
			current = Incident{ID: incidentID(e), Events: []Event{e}}
		}
	}
	return append(incidents, current)
}

func incidentID(e Event) string {
	return e.Timestamp.Format("20060102T150405") + "-" + e.Host
}

// Len returns the number of buffered events.
func (c *Correlator) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.events)
}

// ensure scanner import is indirectly satisfied via callers
var _ = scanner.Result{}
