package debounce_test

import (
	"testing"

	"github.com/user/portwatch/debounce"
	"github.com/user/portwatch/scanner"
)

func makeDiff(opened, closed []scanner.Result) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func result(host string, port int, proto string) scanner.Result {
	return scanner.Result{Host: host, Port: port, Proto: proto, Open: true}
}

func TestBelowThresholdNotConfirmed(t *testing.T) {
	d := debounce.New(3)
	diff := makeDiff([]scanner.Result{result("localhost", 8080, "tcp")}, nil)

	out := d.Evaluate(diff)
	if len(out.Opened) != 0 {
		t.Fatalf("expected 0 confirmed, got %d", len(out.Opened))
	}
	if d.PendingLen() != 1 {
		t.Fatalf("expected 1 pending, got %d", d.PendingLen())
	}
}

func TestAtThresholdConfirmed(t *testing.T) {
	d := debounce.New(2)
	diff := makeDiff([]scanner.Result{result("localhost", 9090, "tcp")}, nil)

	d.Evaluate(diff)
	out := d.Evaluate(diff)

	if len(out.Opened) != 1 {
		t.Fatalf("expected 1 confirmed, got %d", len(out.Opened))
	}
	if d.PendingLen() != 0 {
		t.Fatalf("expected 0 pending after confirmation, got %d", d.PendingLen())
	}
}

func TestThresholdOneConfirmsImmediately(t *testing.T) {
	d := debounce.New(1)
	diff := makeDiff([]scanner.Result{result("localhost", 22, "tcp")}, nil)

	out := d.Evaluate(diff)
	if len(out.Opened) != 1 {
		t.Fatalf("expected immediate confirmation, got %d", len(out.Opened))
	}
}

func TestStaleKeyEvicted(t *testing.T) {
	d := debounce.New(3)
	diff1 := makeDiff([]scanner.Result{result("localhost", 8080, "tcp")}, nil)
	empty := makeDiff(nil, nil)

	d.Evaluate(diff1)
	if d.PendingLen() != 1 {
		t.Fatalf("expected 1 pending")
	}
	d.Evaluate(empty) // key disappears — should be evicted
	if d.PendingLen() != 0 {
		t.Fatalf("expected stale key to be evicted, got %d", d.PendingLen())
	}
}

func TestClosedPortDebounced(t *testing.T) {
	d := debounce.New(2)
	diff := makeDiff(nil, []scanner.Result{result("localhost", 443, "tcp")})

	d.Evaluate(diff)
	out := d.Evaluate(diff)

	if len(out.Closed) != 1 {
		t.Fatalf("expected 1 confirmed closed, got %d", len(out.Closed))
	}
}

func TestFlushClearsAllPending(t *testing.T) {
	d := debounce.New(5)
	diff := makeDiff([]scanner.Result{result("localhost", 8080, "tcp")}, nil)

	d.Evaluate(diff)
	d.Flush()

	if d.PendingLen() != 0 {
		t.Fatalf("expected 0 after flush, got %d", d.PendingLen())
	}
}
