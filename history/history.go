// Package history records scan results over time for trend analysis.
package history

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/scanner"
)

// Entry represents a single historical scan record.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Results   []scanner.Result `json:"results"`
}

// History holds a collection of scan entries.
type History struct {
	Entries []Entry `json:"entries"`
	maxSize int
}

// New creates a new History with a maximum number of entries to retain.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &History{maxSize: maxSize}
}

// Add appends a new scan result set with the current timestamp.
func (h *History) Add(results []scanner.Result) {
	h.Entries = append(h.Entries, Entry{
		Timestamp: time.Now().UTC(),
		Results:   results,
	})
	if len(h.Entries) > h.maxSize {
		h.Entries = h.Entries[len(h.Entries)-h.maxSize:]
	}
}

// Save writes the history to a JSON file.
func (h *History) Save(path string) error {
	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads history from a JSON file. Returns empty history if file missing.
func Load(path string, maxSize int) (*History, error) {
	h := New(maxSize)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return h, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, h); err != nil {
		return nil, err
	}
	return h, nil
}
