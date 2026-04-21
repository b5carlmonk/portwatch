// Package topology builds a map of hosts to their open ports and services,
// providing a structured view of the network surface observed by portwatch.
package topology

import (
	"fmt"
	"sort"
	"sync"

	"github.com/user/portwatch/scanner"
)

// Node represents a single host and the ports observed on it.
type Node struct {
	Host  string
	Ports []scanner.Result
}

// Map holds the current topology keyed by host.
type Map struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

// New returns an empty topology Map.
func New() *Map {
	return &Map{nodes: make(map[string]*Node)}
}

// Update replaces the port list for every host found in results.
func (m *Map) Update(results []scanner.Result) {
	byHost := make(map[string][]scanner.Result)
	for _, r := range results {
		byHost[r.Host] = append(byHost[r.Host], r)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for host, ports := range byHost {
		m.nodes[host] = &Node{Host: host, Ports: ports}
	}
}

// Get returns the Node for a given host and whether it was found.
func (m *Map) Get(host string) (Node, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n, ok := m.nodes[host]
	if !ok {
		return Node{}, false
	}
	return *n, true
}

// Hosts returns all known host addresses in sorted order.
func (m *Map) Hosts() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, 0, len(m.nodes))
	for h := range m.nodes {
		out = append(out, h)
	}
	sort.Strings(out)
	return out
}

// Summary returns a human-readable string describing the topology.
func (m *Map) Summary() string {
	hosts := m.Hosts()
	if len(hosts) == 0 {
		return "topology: no hosts"
	}
	out := fmt.Sprintf("topology: %d host(s)\n", len(hosts))
	for _, h := range hosts {
		n, _ := m.Get(h)
		out += fmt.Sprintf("  %s: %d port(s)\n", h, len(n.Ports))
	}
	return out
}
