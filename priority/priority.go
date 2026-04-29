// Package priority assigns severity levels to port scan changes
// based on configurable rules (port number, protocol, direction).
package priority

import (
	"fmt"
	"strings"

	"github.com/user/portwatch/scanner"
)

// Level represents the severity of a change event.
type Level int

const (
	Low Level = iota
	Medium
	High
	Critical
)

func (l Level) String() string {
	switch l {
	case Critical:
		return "CRITICAL"
	case High:
		return "HIGH"
	case Medium:
		return "MEDIUM"
	default:
		return "LOW"
	}
}

// Rule maps a port/protocol pair to a Level.
type Rule struct {
	Port     int
	Proto    string // "tcp", "udp", or "" for any
	Level    Level
}

// Scorer assigns priority levels to scan results.
type Scorer struct {
	rules    []Rule
	fallback Level
}

// New returns a Scorer with the given rules and a fallback level for
// ports that match no rule.
func New(rules []Rule, fallback Level) *Scorer {
	return &Scorer{rules: rules, fallback: fallback}
}

// Score returns the Level for the given result.
func (s *Scorer) Score(r scanner.Result) Level {
	for _, rule := range s.rules {
		protoMatch := rule.Proto == "" || strings.EqualFold(rule.Proto, r.Proto)
		if rule.Port == r.Port && protoMatch {
			return rule.Level
		}
	}
	return s.fallback
}

// ScoreAll returns a map of "port/proto" -> Level for every result.
func (s *Scorer) ScoreAll(results []scanner.Result) map[string]Level {
	out := make(map[string]Level, len(results))
	for _, r := range results {
		k := fmt.Sprintf("%d/%s", r.Port, strings.ToLower(r.Proto))
		out[k] = s.Score(r)
	}
	return out
}
