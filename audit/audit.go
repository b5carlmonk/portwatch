// Package audit records a timestamped log of every scan cycle result.
package audit

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/scanner"
)

// Entry is a single audit record.
type Entry struct {
	Time    time.Time        `json:"time"`
	Host    string           `json:"host"`
	Results []scanner.Result `json:"results"`
}

// Log holds an ordered list of audit entries.
type Log struct {
	mu      sync.Mutex
	entries []Entry
	path    string
}

// New creates a new Log that persists to path.
func New(path string) *Log {
	return &Log{path: path}
}

// Record appends a new entry for the current time.
func (l *Log) Record(host string, results []scanner.Result) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, Entry{
		Time:    time.Now().UTC(),
		Host:    host,
		Results: results,
	})
}

// Entries returns a copy of all recorded entries.
func (l *Log) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// Save writes the log to disk as JSON.
func (l *Log) Save() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.Create(l.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(l.entries)
}

// Load reads a previously saved log from disk.
func (l *Log) Load() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	f, err := os.Open(l.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&l.entries)
}
