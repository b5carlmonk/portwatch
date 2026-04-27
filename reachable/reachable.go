// Package reachable checks whether scanned hosts are network-reachable
// before a full port scan is attempted, avoiding wasted scan cycles.
package reachable

import (
	"context"
	"net"
	"sync"
	"time"
)

// Result holds the outcome of a reachability probe for a single host.
type Result struct {
	Host      string
	Reachable bool
	Latency   time.Duration
	Err       error
}

// Checker probes hosts for reachability via a lightweight TCP dial.
type Checker struct {
	timeout    time.Duration
	probePort  string
	mu         sync.Mutex
	cache      map[string]Result
	cacheTTL   time.Duration
	cachedAt   map[string]time.Time
}

// New returns a Checker with the given dial timeout and probe port.
// Use probePort "80" or "443" for general reachability checks.
func New(timeout time.Duration, probePort string) *Checker {
	return &Checker{
		timeout:   timeout,
		probePort: probePort,
		cache:     make(map[string]Result),
		cacheTTL:  30 * time.Second,
		cachedAt:  make(map[string]time.Time),
	}
}

// Probe checks whether host is reachable. Results are cached for the TTL
// duration to avoid redundant network calls within a single scan cycle.
func (c *Checker) Probe(ctx context.Context, host string) Result {
	c.mu.Lock()
	if r, ok := c.cache[host]; ok {
		if time.Since(c.cachedAt[host]) < c.cacheTTL {
			c.mu.Unlock()
			return r
		}
	}
	c.mu.Unlock()

	addr := net.JoinHostPort(host, c.probePort)
	start := time.Now()
	conn, err := (&net.Dialer{Timeout: c.timeout}).DialContext(ctx, "tcp", addr)
	latency := time.Since(start)

	r := Result{Host: host, Latency: latency}
	if err != nil {
		r.Err = err
		r.Reachable = false
	} else {
		conn.Close()
		r.Reachable = true
	}

	c.mu.Lock()
	c.cache[host] = r
	c.cachedAt[host] = time.Now()
	c.mu.Unlock()

	return r
}

// ProbeAll probes each host concurrently and returns all results.
func (c *Checker) ProbeAll(ctx context.Context, hosts []string) []Result {
	results := make([]Result, len(hosts))
	var wg sync.WaitGroup
	for i, h := range hosts {
		wg.Add(1)
		go func(idx int, host string) {
			defer wg.Done()
			results[idx] = c.Probe(ctx, host)
		}(i, h)
	}
	wg.Wait()
	return results
}

// Flush clears the result cache.
func (c *Checker) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]Result)
	c.cachedAt = make(map[string]time.Time)
}
