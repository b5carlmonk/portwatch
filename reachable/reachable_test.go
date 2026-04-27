package reachable_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/reachable"
)

func startTCPServer(t *testing.T) (string, string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	host, port, _ := net.SplitHostPort(ln.Addr().String())
	return host, port
}

func TestProbeReachableHost(t *testing.T) {
	host, port := startTCPServer(t)
	c := reachable.New(2*time.Second, port)
	r := c.Probe(context.Background(), host)
	if !r.Reachable {
		t.Fatalf("expected host to be reachable, got err: %v", r.Err)
	}
	if r.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestProbeUnreachableHost(t *testing.T) {
	c := reachable.New(200*time.Millisecond, "19999")
	r := c.Probe(context.Background(), "127.0.0.1")
	if r.Reachable {
		t.Fatal("expected host to be unreachable")
	}
	if r.Err == nil {
		t.Error("expected non-nil error for unreachable host")
	}
}

func TestProbeCachesResult(t *testing.T) {
	host, port := startTCPServer(t)
	c := reachable.New(2*time.Second, port)
	r1 := c.Probe(context.Background(), host)
	r2 := c.Probe(context.Background(), host)
	if r1.Reachable != r2.Reachable {
		t.Error("cached result should match original")
	}
}

func TestFlushClearsCache(t *testing.T) {
	host, port := startTCPServer(t)
	c := reachable.New(2*time.Second, port)
	c.Probe(context.Background(), host)
	c.Flush()
	// After flush a fresh probe should still succeed
	r := c.Probe(context.Background(), host)
	if !r.Reachable {
		t.Fatalf("expected reachable after flush, got: %v", r.Err)
	}
}

func TestProbeAllReturnsAllResults(t *testing.T) {
	host, port := startTCPServer(t)
	c := reachable.New(2*time.Second, port)
	results := c.ProbeAll(context.Background(), []string{host, host})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Reachable {
			t.Errorf("expected reachable, got err: %v", r.Err)
		}
	}
}

func TestProbeAllEmptyInput(t *testing.T) {
	c := reachable.New(time.Second, "80")
	results := c.ProbeAll(context.Background(), nil)
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}
