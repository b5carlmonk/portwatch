// Package sampler provides periodic port scan sampling with configurable
// intervals, storing named samples for later comparison or trend analysis.
package sampler

import (
	"fmt"
	"sync"
	"time"

	"github.com/user/portwatch/scanner"
)

// Sample holds a named snapshot of scan results captured at a point in time.
type Sample struct {
	Name      string
	CapturedAt time.Time
	Results   []scanner.Result
}

// Sampler captures and stores named scan samples.
type Sampler struct {
	mu      sync.RWMutex
	samples map[string]Sample
	scan    func(host string, ports []int) ([]scanner.Result, error)
}

// New returns a Sampler that uses the provided scan function to collect results.
func New(scan func(host string, ports []int) ([]scanner.Result, error)) *Sampler {
	return &Sampler{
		samples: make(map[string]Sample),
		scan:    scan,
	}
}

// Capture runs a scan against host and ports, storing the results under name.
func (s *Sampler) Capture(name, host string, ports []int) error {
	if name == "" {
		return fmt.Errorf("sampler: name must not be empty")
	}
	results, err := s.scan(host, ports)
	if err != nil {
		return fmt.Errorf("sampler: scan failed: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples[name] = Sample{
		Name:       name,
		CapturedAt: time.Now(),
		Results:    results,
	}
	return nil
}

// Get returns the sample stored under name, or an error if not found.
func (s *Sampler) Get(name string) (Sample, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sample, ok := s.samples[name]
	if !ok {
		return Sample{}, fmt.Errorf("sampler: no sample named %q", name)
	}
	return sample, nil
}

// Delete removes the sample stored under name.
func (s *Sampler) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.samples, name)
}

// Names returns the names of all stored samples.
func (s *Sampler) Names() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.samples))
	for k := range s.samples {
		names = append(names, k)
	}
	return names
}
