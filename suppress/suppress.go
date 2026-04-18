// Package suppress provides a flap-detection mechanism that prevents
// alerts from firing until a change has been observed N consecutive times.
package suppress

import "sync"

// Suppressor tracks consecutive occurrences of a keyed event and only
// reports it as "confirmed" once the threshold is reached.
type Suppressor struct {
	mu        sync.Mutex
	threshold int
	counts    map[string]int
}

// New returns a Suppressor that confirms an event after threshold hits.
func New(threshold int) *Suppressor {
	if threshold < 1 {
		threshold = 1
	}
	return &Suppressor{
		threshold: threshold,
		counts:    make(map[string]int),
	}
}

// Record increments the counter for key and returns true when the
// threshold has been reached exactly (i.e. the first confirmed firing).
func (s *Suppressor) Record(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[key]++
	return s.counts[key] == s.threshold
}

// Reset clears the counter for key (e.g. when a change disappears).
func (s *Suppressor) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.counts, key)
}

// Count returns the current consecutive count for key.
func (s *Suppressor) Count(key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counts[key]
}

// Flush clears all tracked counters.
func (s *Suppressor) Flush() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[string]int)
}
