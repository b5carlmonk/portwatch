// Package filter provides port filtering utilities for portwatch.
package filter

import "github.com/user/portwatch/scanner"

// Rule defines a single filter rule.
type Rule struct {
	Port     int
	Protocol string // "tcp" or "udp", empty means any
}

// Filter holds a set of rules to include or exclude ports.
type Filter struct {
	Include []Rule
	Exclude []Rule
}

// New returns a Filter with the given include and exclude rules.
func New(include, exclude []Rule) *Filter {
	return &Filter{Include: include, Exclude: exclude}
}

// Apply returns only the scan results that pass the filter.
// If Include is non-empty, only matching results are kept.
// Results matching any Exclude rule are removed.
func (f *Filter) Apply(results []scanner.Result) []scanner.Result {
	var out []scanner.Result
	for _, r := range results {
		if len(f.Include) > 0 && !f.matchAny(r, f.Include) {
			continue
		}
		if f.matchAny(r, f.Exclude) {
			continue
		}
		out = append(out, r)
	}
	return out
}

func (f *Filter) matchAny(r scanner.Result, rules []Rule) bool {
	for _, rule := range rules {
		if rule.Port != 0 && rule.Port != r.Port {
			continue
		}
		if rule.Protocol != "" && rule.Protocol != r.Protocol {
			continue
		}
		return true
	}
	return false
}
