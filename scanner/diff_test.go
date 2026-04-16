package scanner_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/scanner"
)

func makeResult(host string, ports []scanner.PortState) *scanner.ScanResult {
	return &scanner.ScanResult{
		Host:      host,
		Timestamp: time.Now(),
		Ports:     ports,
	}
}

func TestDiffDetectsOpenedPort(t *testing.T) {
	prev := makeResult("localhost", []scanner.PortState{
		{Port: 80, Protocol: "tcp", Open: false},
	})
	curr := makeResult("localhost", []scanner.PortState{
		{Port: 80, Protocol: "tcp", Open: true, Service: "http"},
	})

	changes := scanner.Diff(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Change != scanner.ChangeOpened {
		t.Errorf("expected ChangeOpened, got %s", changes[0].Change)
	}
}

func TestDiffDetectsClosedPort(t *testing.T) {
	prev := makeResult("localhost", []scanner.PortState{
		{Port: 22, Protocol: "tcp", Open: true, Service: "ssh"},
	})
	curr := makeResult("localhost", []scanner.PortState{
		{Port: 22, Protocol: "tcp", Open: false},
	})

	changes := scanner.Diff(prev, curr)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Change != scanner.ChangeClosed {
		t.Errorf("expected ChangeClosed, got %s", changes[0].Change)
	}
}

func TestDiffNoChanges(t *testing.T) {
	prev := makeResult("localhost", []scanner.PortState{
		{Port: 443, Protocol: "tcp", Open: true, Service: "https"},
	})
	curr := makeResult("localhost", []scanner.PortState{
		{Port: 443, Protocol: "tcp", Open: true, Service: "https"},
	})

	changes := scanner.Diff(prev, curr)
	if len(changes) != 0 {
		t.Errorf("expected no changes, got %d", len(changes))
	}
}
