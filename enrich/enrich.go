// Package enrich attaches extra metadata to scan results.
package enrich

import (
	"fmt"
	"time"

	"github.com/user/portwatch/scanner"
)

// Meta holds additional metadata attached to a scan result.
type Meta struct {
	Label     string
	ScannedAt time.Time
	Extra     map[string]string
}

// Enriched wraps a scanner.Result with additional metadata.
type Enriched struct {
	scanner.Result
	Meta Meta
}

// Enricher applies metadata to scan results.
type Enricher struct {
	providers []Provider
}

// Provider is a function that returns key/value metadata for a result.
type Provider func(r scanner.Result) (key, value string)

// New creates a new Enricher with the given providers.
func New(providers ...Provider) *Enricher {
	return &Enricher{providers: providers}
}

// Enrich attaches metadata from all providers to each result.
func (e *Enricher) Enrich(results []scanner.Result) []Enriched {
	out := make([]Enriched, 0, len(results))
	for _, r := range results {
		meta := Meta{
			ScannedAt: time.Now(),
			Extra:     make(map[string]string),
		}
		for _, p := range e.providers {
			k, v := p(r)
			if k != "" {
				meta.Extra[k] = v
			}
		}
		out = append(out, Enriched{Result: r, Meta: meta})
	}
	return out
}

// PortLabel is a built-in provider that labels results by port number.
func PortLabel(r scanner.Result) (string, string) {
	return "port_label", fmt.Sprintf("port-%d", r.Port)
}

// ProtoLabel is a built-in provider that labels results by protocol.
func ProtoLabel(r scanner.Result) (string, string) {
	return "proto_label", r.Proto
}
