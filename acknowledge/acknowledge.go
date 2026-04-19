// Package acknowledge provides a simple store for suppressing alerts
// on known/acknowledged ports so operators can silence expected changes.
package acknowledge

import (
	"encoding/json"
	"os"
	"sync"
)

// Key identifies an acknowledged port event.
type Key struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

// Store holds acknowledged keys.
type Store struct {
	mu   sync.RWMutex
	keys map[Key]struct{}
	path string
}

// New returns an empty Store backed by the given file path.
func New(path string) *Store {
	return &Store{path: path, keys: make(map[Key]struct{})}
}

// Acknowledge marks a key as acknowledged.
func (s *Store) Acknowledge(k Key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.keys[k] = struct{}{}
}

// IsAcknowledged reports whether the key has been acknowledged.
func (s *Store) IsAcknowledged(k Key) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.keys[k]
	return ok
}

// Revoke removes an acknowledgement.
func (s *Store) Revoke(k Key) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.keys, k)
}

// Save persists the store to disk.
func (s *Store) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]Key, 0, len(s.keys))
	for k := range s.keys {
		list = append(list, k)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Load reads acknowledged keys from disk.
func (s *Store) Load() error {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	var list []Key
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, k := range list {
		s.keys[k] = struct{}{}
	}
	return nil
}
