package resolve_test

import (
	"fmt"
	"testing"

	"github.com/user/portwatch/resolve"
	"github.com/user/portwatch/scanner"
)

func fakeLookup(hits map[string]string) func(string) ([]string, error) {
	return func(ip string) ([]string, error) {
		if name, ok := hits[ip]; ok {
			return []string{name}, nil
		}
		return nil, fmt.Errorf("not found")
	}
}

func TestLookupReturnsHostname(t *testing.T) {
	r := resolve.NewWithLookup(fakeLookup(map[string]string{"1.2.3.4": "example.com."}))
	if got := r.Lookup("1.2.3.4"); got != "example.com." {
		t.Fatalf("expected example.com. got %s", got)
	}
}

func TestLookupFallsBackToIP(t *testing.T) {
	r := resolve.NewWithLookup(fakeLookup(map[string]string{}))
	if got := r.Lookup("9.9.9.9"); got != "9.9.9.9" {
		t.Fatalf("expected IP fallback got %s", got)
	}
}

func TestLookupCaches(t *testing.T) {
	calls := 0
	r := resolve.NewWithLookup(func(ip string) ([]string, error) {
		calls++
		return []string{"cached.host."}, nil
	})
	r.Lookup("1.1.1.1")
	r.Lookup("1.1.1.1")
	if calls != 1 {
		t.Fatalf("expected 1 DNS call, got %d", calls)
	}
}

func TestEnrichUpdatesHost(t *testing.T) {
	r := resolve.NewWithLookup(fakeLookup(map[string]string{"10.0.0.1": "internal.local."}))
	results := []scanner.Result{
		{Host: "10.0.0.1", Port: 80, Proto: "tcp", Open: true},
	}
	enriched := r.Enrich(results)
	if enriched[0].Host != "internal.local." {
		t.Fatalf("unexpected host %s", enriched[0].Host)
	}
}

func TestEnrichPreservesOtherFields(t *testing.T) {
	r := resolve.NewWithLookup(fakeLookup(map[string]string{}))
	results := []scanner.Result{
		{Host: "5.5.5.5", Port: 443, Proto: "tcp", Open: true},
	}
	enriched := r.Enrich(results)
	if enriched[0].Port != 443 || enriched[0].Proto != "tcp" {
		t.Fatal("fields mutated unexpectedly")
	}
}
