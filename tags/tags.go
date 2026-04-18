// Package tags provides labeling support for scan results.
package tags

import (
	"strings"

	"github.com/user/portwatch/scanner"
)

// Tag represents a key=value label attached to a scan result.
type Tag struct {
	Key   string
	Value string
}

// Tagger applies tags to scan results based on configured rules.
type Tagger struct {
	rules []rule
}

type rule struct {
	port  int
	proto string
	tags  []Tag
}

// New returns a Tagger with no rules.
func New() *Tagger {
	return &Tagger{}
}

// AddRule registers tags to apply when port and proto match.
// An empty proto matches any protocol.
func (t *Tagger) AddRule(port int, proto string, tags ...Tag) {
	t.rules = append(t.rules, rule{port: port, proto: strings.ToLower(proto), tags: tags})
}

// Apply returns a map from result key to matching tags for each result.
func (t *Tagger) Apply(results []scanner.Result) map[string][]Tag {
	out := make(map[string][]Tag)
	for _, r := range results {
		key := resultKey(r)
		for _, ru := range t.rules {
			if ru.port != r.Port {
				continue
			}
			if ru.proto != "" && !strings.EqualFold(ru.proto, r.Proto) {
				continue
			}
			out[key] = append(out[key], ru.tags...)
		}
	}
	return out
}

func resultKey(r scanner.Result) string {
	return r.Host + ":" + strings.ToLower(r.Proto) + ":" + itoa(r.Port)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
