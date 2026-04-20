package correlate

import (
	"testing"
	"time"
)

func TestAddIncreasesLen(t *testing.T) {
	c := New(time.Second, nil)
	c.Add("host1", 80, "tcp", "opened")
	c.Add("host2", 443, "tcp", "opened")
	if c.Len() != 2 {
		t.Fatalf("expected 2 events, got %d", c.Len())
	}
}

func TestFlushClearsBuffer(t *testing.T) {
	c := New(time.Second, nil)
	c.Add("host1", 22, "tcp", "opened")
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected 0 events after flush, got %d", c.Len())
	}
}

func TestFlushEmptyDoesNotPanic(t *testing.T) {
	c := New(time.Second, nil)
	c.Flush() // should not panic
}

func TestEventsWithinWindowGroupedTogether(t *testing.T) {
	var incidents []Incident
	c := New(time.Minute, func(inc Incident) {
		incidents = append(incidents, inc)
	})
	now := time.Now()
	c.mu.Lock()
	c.events = []Event{
		{Host: "h1", Port: 80, Proto: "tcp", State: "opened", Timestamp: now},
		{Host: "h2", Port: 80, Proto: "tcp", State: "opened", Timestamp: now.Add(5 * time.Second)},
		{Host: "h3", Port: 80, Proto: "tcp", State: "opened", Timestamp: now.Add(10 * time.Second)},
	}
	c.mu.Unlock()
	c.Flush()
	if len(incidents) != 1 {
		t.Fatalf("expected 1 incident, got %d", len(incidents))
	}
	if len(incidents[0].Events) != 3 {
		t.Fatalf("expected 3 events in incident, got %d", len(incidents[0].Events))
	}
}

func TestEventsOutsideWindowSplitIntoMultipleIncidents(t *testing.T) {
	var incidents []Incident
	c := New(5*time.Second, func(inc Incident) {
		incidents = append(incidents, inc)
	})
	now := time.Now()
	c.mu.Lock()
	c.events = []Event{
		{Host: "h1", Port: 22, Proto: "tcp", State: "opened", Timestamp: now},
		{Host: "h2", Port: 22, Proto: "tcp", State: "opened", Timestamp: now.Add(30 * time.Second)},
	}
	c.mu.Unlock()
	c.Flush()
	if len(incidents) != 2 {
		t.Fatalf("expected 2 incidents, got %d", len(incidents))
	}
}

func TestIncidentIDIsNonEmpty(t *testing.T) {
	e := Event{Host: "myhost", Port: 443, Proto: "tcp", State: "closed", Timestamp: time.Now()}
	id := incidentID(e)
	if id == "" {
		t.Fatal("expected non-empty incident ID")
	}
}

func TestCallbackReceivesCorrectHost(t *testing.T) {
	var got []string
	c := New(time.Minute, func(inc Incident) {
		for _, e := range inc.Events {
			got = append(got, e.Host)
		}
	})
	c.Add("alpha", 8080, "tcp", "opened")
	c.Flush()
	if len(got) != 1 || got[0] != "alpha" {
		t.Fatalf("unexpected hosts: %v", got)
	}
}
