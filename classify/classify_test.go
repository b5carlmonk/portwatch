package classify_test

import (
	"testing"

	"github.com/user/portwatch/classify"
	"github.com/user/portwatch/scanner"
)

func makeResults(ports []int, proto string) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Host: "127.0.0.1", Port: p, Proto: proto, Open: true}
	}
	return out
}

func TestLabelKnownPort(t *testing.T) {
	c := classify.New(nil)
	if got := c.Label(22, "tcp"); got != "SSH" {
		t.Fatalf("expected SSH, got %s", got)
	}
}

func TestLabelUnknownPort(t *testing.T) {
	c := classify.New(nil)
	if got := c.Label(9999, "tcp"); got != "unknown" {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestLabelCustomOverridesDefault(t *testing.T) {
	c := classify.New(map[string]string{"22/tcp": "MySSH"})
	if got := c.Label(22, "tcp"); got != "MySSH" {
		t.Fatalf("expected MySSH, got %s", got)
	}
}

func TestLabelCustomNewEntry(t *testing.T) {
	c := classify.New(map[string]string{"9200/tcp": "Elasticsearch"})
	if got := c.Label(9200, "tcp"); got != "Elasticsearch" {
		t.Fatalf("expected Elasticsearch, got %s", got)
	}
}

func TestEnrichSetsServiceField(t *testing.T) {
	c := classify.New(nil)
	results := makeResults([]int{80, 443, 9999}, "tcp")
	enriched := c.Enrich(results)

	expected := map[int]string{80: "HTTP", 443: "HTTPS", 9999: "unknown"}
	for _, r := range enriched {
		want, ok := expected[r.Port]
		if !ok {
			t.Fatalf("unexpected port %d", r.Port)
		}
		if r.Service != want {
			t.Errorf("port %d: expected service %q, got %q", r.Port, want, r.Service)
		}
	}
}

func TestEnrichEmptySlice(t *testing.T) {
	c := classify.New(nil)
	out := c.Enrich([]scanner.Result{})
	if len(out) != 0 {
		t.Fatalf("expected empty slice, got %d entries", len(out))
	}
}

func TestEnrichUDPPort(t *testing.T) {
	c := classify.New(nil)
	results := []scanner.Result{{Host: "127.0.0.1", Port: 53, Proto: "udp", Open: true}}
	enriched := c.Enrich(results)
	if enriched[0].Service != "DNS" {
		t.Errorf("expected DNS, got %s", enriched[0].Service)
	}
}

// TestEnrichPreservesOriginalFields verifies that Enrich does not mutate fields
// other than Service on the returned results.
func TestEnrichPreservesOriginalFields(t *testing.T) {
	c := classify.New(nil)
	input := []scanner.Result{
		{Host: "10.0.0.1", Port: 22, Proto: "tcp", Open: true},
	}
	enriched := c.Enrich(input)
	if len(enriched) != 1 {
		t.Fatalf("expected 1 result, got %d", len(enriched))
	}
	r := enriched[0]
	if r.Host != "10.0.0.1" {
		t.Errorf("Host changed: got %q", r.Host)
	}
	if r.Port != 22 {
		t.Errorf("Port changed: got %d", r.Port)
	}
	if r.Proto != "tcp" {
		t.Errorf("Proto changed: got %q", r.Proto)
	}
	if !r.Open {
		t.Errorf("Open changed: got false")
	}
}
