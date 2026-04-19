package groupby_test

import (
	"testing"

	"github.com/user/portwatch/groupby"
	"github.com/user/portwatch/scanner"
)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "localhost", Port: 80, Protocol: "tcp", Open: true},
		{Host: "localhost", Port: 443, Protocol: "tcp", Open: true},
		{Host: "remotehost", Port: 80, Protocol: "udp", Open: true},
		{Host: "remotehost", Port: 53, Protocol: "udp", Open: true},
	}
}

func TestGroupByPort(t *testing.T) {
	g := groupby.New(groupby.ByPort)
	groups := g.Group(makeResults())
	if len(groups["80"]) != 2 {
		t.Errorf("expected 2 results for port 80, got %d", len(groups["80"]))
	}
	if len(groups["443"]) != 1 {
		t.Errorf("expected 1 result for port 443, got %d", len(groups["443"]))
	}
}

func TestGroupByProtocol(t *testing.T) {
	g := groupby.New(groupby.ByProtocol)
	groups := g.Group(makeResults())
	if len(groups["tcp"]) != 2 {
		t.Errorf("expected 2 tcp results, got %d", len(groups["tcp"]))
	}
	if len(groups["udp"]) != 2 {
		t.Errorf("expected 2 udp results, got %d", len(groups["udp"]))
	}
}

func TestGroupByHost(t *testing.T) {
	g := groupby.New(groupby.ByHost)
	groups := g.Group(makeResults())
	if len(groups["localhost"]) != 2 {
		t.Errorf("expected 2 localhost results, got %d", len(groups["localhost"]))
	}
	if len(groups["remotehost"]) != 2 {
		t.Errorf("expected 2 remotehost results, got %d", len(groups["remotehost"]))
	}
}

func TestGroupEmptyInput(t *testing.T) {
	g := groupby.New(groupby.ByPort)
	groups := g.Group([]scanner.Result{})
	if len(groups) != 0 {
		t.Errorf("expected empty map, got %d keys", len(groups))
	}
}
