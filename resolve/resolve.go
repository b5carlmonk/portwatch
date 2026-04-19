// Package resolve performs reverse DNS lookups for scan results.
package resolve

import (
	"net"
	"sync"

	"github.com/user/portwatch/scanner"
)

// Resolver maps IPs to hostnames.
type Resolver struct {
	mu    sync.Mutex
	cache map[string]string
	lookup func(string) ([]string, error)
}

// New returns a Resolver using the system DNS.
func New() *Resolver {
	return &Resolver{
		cache:  make(map[string]string),
		lookup: net.LookupAddr,
	}
}

// NewWithLookup returns a Resolver using a custom lookup function (useful for tests).
func NewWithLookup(fn func(string) ([]string, error)) *Resolver {
	return &Resolver{
		cache:  make(map[string]string),
		lookup: fn,
	}
}

// Lookup returns the hostname for an IP, using a cache to avoid repeated lookups.
func (r *Resolver) Lookup(ip string) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if v, ok := r.cache[ip]; ok {
		return v
	}
	names, err := r.lookup(ip)
	if err != nil || len(names) == 0 {
		r.cache[ip] = ip
		return ip
	}
	host := names[0]
	r.cache[ip] = host
	return host
}

// Enrich annotates each ScanResult's Host field with a resolved hostname when
// the host is a bare IP address.
func (r *Resolver) Enrich(results []scanner.Result) []scanner.Result {
	out := make([]scanner.Result, len(results))
	for i, res := range results {
		res.Host = r.Lookup(res.Host)
		out[i] = res
	}
	return out
}
