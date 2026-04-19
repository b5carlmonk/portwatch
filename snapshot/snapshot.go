// Package snapshot captures and compares point-in-time port scan results.
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/portwatch/scanner"
)

// Entry holds a labeled snapshot of scan results.
type Entry struct {
	Label     string           `json:"label"`
	CreatedAt time.Time        `json:"created_at"`
	Results   []scanner.Result `json:"results"`
}

// Store manages named snapshots persisted to disk.
type Store struct {
	path string
	data map[string]Entry
}

// New returns a Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path, data: make(map[string]Entry)}
}

// Add stores a snapshot under the given label.
func (s *Store) Add(label string, results []scanner.Result) error {
	if label == "" {
		return fmt.Errorf("snapshot: label must not be empty")
	}
	s.data[label] = Entry{
		Label:     label,
		CreatedAt: time.Now().UTC(),
		Results:   results,
	}
	return s.save()
}

// Get returns the snapshot for the given label.
func (s *Store) Get(label string) (Entry, bool) {
	e, ok := s.data[label]
	return e, ok
}

// Delete removes a snapshot by label.
func (s *Store) Delete(label string) error {
	delete(s.data, label)
	return s.save()
}

// Labels returns all stored snapshot labels.
func (s *Store) Labels() []string {
	out := make([]string, 0, len(s.data))
	for k := range s.data {
		out = append(out, k)
	}
	return out
}

func (s *Store) save() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}

// Load reads persisted snapshots from disk.
func (s *Store) Load() error {
	b, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.data)
}
