// Package trend tracks port open/close frequency over time.
package trend

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry records a single observation for a port key.
type Entry struct {
	Key       string    `json:"key"`
	OpenCount int       `json:"open_count"`
	CloseCount int      `json:"close_count"`
	LastSeen  time.Time `json:"last_seen"`
}

// Tracker holds trend data keyed by "host:port/proto".
type Tracker struct {
	mu      sync.Mutex
	entries map[string]*Entry
}

// New returns an empty Tracker.
func New() *Tracker {
	return &Tracker{entries: make(map[string]*Entry)}
}

// RecordOpen increments the open counter for key.
func (t *Tracker) RecordOpen(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e := t.get(key)
	e.OpenCount++
	e.LastSeen = time.Now()
}

// RecordClose increments the close counter for key.
func (t *Tracker) RecordClose(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e := t.get(key)
	e.CloseCount++
	e.LastSeen = time.Now()
}

// Get returns a copy of the entry for key, and whether it exists.
func (t *Tracker) Get(key string) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[key]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// All returns a snapshot of all entries.
func (t *Tracker) All() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

func (t *Tracker) get(key string) *Entry {
	if e, ok := t.entries[key]; ok {
		return e
	}
	e := &Entry{Key: key}
	t.entries[key] = e
	return e
}

// Save persists the tracker state to path.
func (t *Tracker) Save(path string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(t.entries)
}

// Load restores tracker state from path.
func Load(path string) (*Tracker, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return New(), nil
		}
		return nil, err
	}
	defer f.Close()
	t := New()
	if err := json.NewDecoder(f).Decode(&t.entries); err != nil {
		return nil, err
	}
	return t, nil
}
