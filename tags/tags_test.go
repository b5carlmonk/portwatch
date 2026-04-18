package tags_test

import (
	"testing"

	"github.com/user/portwatch/scanner"
	"github.com/user/portwatch/tags"
)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "localhost", Port: 80, Proto: "tcp", Open: true},
		{Host: "localhost", Port: 443, Proto: "tcp", Open: true},
		{Host: "localhost", Port: 53, Proto: "udp", Open: true},
	}
}

func TestNoRulesReturnsEmptyMap(t *testing.T) {
	tr := tags.New()
	out := tr.Apply(makeResults())
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(out))
	}
}

func TestMatchByPortAndProto(t *testing.T) {
	tr := tags.New()
	tr.AddRule(80, "tcp", tags.Tag{Key: "service", Value: "http"})
	out := tr.Apply(makeResults())
	key := "localhost:tcp:80"
	if len(out[key]) != 1 {
		t.Fatalf("expected 1 tag for %s, got %d", key, len(out[key]))
	}
	if out[key][0].Value != "http" {
		t.Errorf("expected value http, got %s", out[key][0].Value)
	}
}

func TestMatchByPortAnyProto(t *testing.T) {
	tr := tags.New()
	tr.AddRule(53, "", tags.Tag{Key: "role", Value: "dns"})
	out := tr.Apply(makeResults())
	key := "localhost:udp:53"
	if len(out[key]) != 1 {
		t.Fatalf("expected 1 tag for %s, got %d", key, len(out[key]))
	}
}

func TestNoMatchOnWrongProto(t *testing.T) {
	tr := tags.New()
	tr.AddRule(53, "tcp", tags.Tag{Key: "role", Value: "dns"})
	out := tr.Apply(makeResults())
	if len(out) != 0 {
		t.Errorf("expected no matches, got %d", len(out))
	}
}

func TestMultipleTagsOnSamePort(t *testing.T) {
	tr := tags.New()
	tr.AddRule(443, "tcp",
		tags.Tag{Key: "service", Value: "https"},
		tags.Tag{Key: "secure", Value: "true"},
	)
	out := tr.Apply(makeResults())
	key := "localhost:tcp:443"
	if len(out[key]) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(out[key]))
	}
}
