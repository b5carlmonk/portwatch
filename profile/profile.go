// Package profile manages named scan profiles (sets of ports/config to scan).
package profile

import (
	"encoding/json"
	"errors"
	"os"
)

// Profile represents a named scanning profile.
type Profile struct {
	Name     string   `json:"name"`
	Hosts    []string `json:"hosts"`
	Ports    []int    `json:"ports"`
	Protocol string   `json:"protocol"`
}

// Store holds multiple named profiles.
type Store struct {
	Profiles map[string]Profile `json:"profiles"`
}

// New returns an empty Store.
func New() *Store {
	return &Store{Profiles: make(map[string]Profile)}
}

// Add inserts or replaces a profile by name.
func (s *Store) Add(p Profile) error {
	if p.Name == "" {
		return errors.New("profile name must not be empty")
	}
	s.Profiles[p.Name] = p
	return nil
}

// Get retrieves a profile by name.
func (s *Store) Get(name string) (Profile, bool) {
	p, ok := s.Profiles[name]
	return p, ok
}

// Remove deletes a profile by name.
func (s *Store) Remove(name string) {
	delete(s.Profiles, name)
}

// Save persists the store to a JSON file.
func (s *Store) Save(path string) error {
	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

// Load reads a store from a JSON file.
func Load(path string) (*Store, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return New(), nil
		}
		return nil, err
	}
	s := New()
	if err := json.Unmarshal(b, s); err != nil {
		return nil, err
	}
	return s, nil
}
