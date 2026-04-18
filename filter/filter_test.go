package filter_test

import (
	"testing"

	"github.com/user/portwatch/filter"
	"github.com/user/portwatch/scanner"
)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "localhost", Port: 22, Protocol: "tcp", Open: true},
		{Host: "localhost", Port: 80, Protocol: "tcp", Open: true},
		{Host: "localhost", Port: 443, Protocol: "tcp", Open: true},
		{Host: "localhost", Port: 53, Protocol: "udp", Open: true},
	}
}

func TestNoRulesReturnsAll(t *testing.T) {
	f := filter.New(nil, nil)
	got := f.Apply(makeResults())
	if len(got) != 4 {
		t.Fatalf("expected 4 results, got %d", len(got))
	}
}

func TestIncludeByPort(t *testing.T) {
	f := filter.New([]filter.Rule{{Port: 80}}, nil)
	got := f.Apply(makeResults())
	if len(got) != 1 || got[0].Port != 80 {
		t.Fatalf("expected port 80, got %v", got)
	}
}

func TestExcludeByPort(t *testing.T) {
	f := filter.New(nil, []filter.Rule{{Port: 22}})
	got := f.Apply(makeResults())
	for _, r := range got {
		if r.Port == 22 {
			t.Fatal("port 22 should have been excluded")
		}
	}
	if len(got) != 3 {
		t.Fatalf("expected 3 results, got %d", len(got))
	}
}

func TestIncludeByProtocol(t *testing.T) {
	f := filter.New([]filter.Rule{{Protocol: "udp"}}, nil)
	got := f.Apply(makeResults())
	if len(got) != 1 || got[0].Protocol != "udp" {
		t.Fatalf("expected only udp result, got %v", got)
	}
}

func TestExcludeByProtocol(t *testing.T) {
	f := filter.New(nil, []filter.Rule{{Protocol: "udp"}})
	got := f.Apply(makeResults())
	for _, r := range got {
		if r.Protocol == "udp" {
			t.Fatal("udp results should have been excluded")
		}
	}
}

func TestIncludeAndExcludeCombined(t *testing.T) {
	// Include tcp, but exclude port 443
	f := filter.New(
		[]filter.Rule{{Protocol: "tcp"}},
		[]filter.Rule{{Port: 443}},
	)
	got := f.Apply(makeResults())
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
	for _, r := range got {
		if r.Port == 443 {
			t.Fatal("port 443 should have been excluded")
		}
	}
}
