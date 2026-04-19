package healthcheck

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func startTCPServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ := strconv.Atoi(portStr)
	return port, func() { ln.Close() }
}

func TestProbeAlive(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	c := New(time.Second)
	r := c.Probe("127.0.0.1", port, "tcp")
	if !r.Alive {
		t.Errorf("expected alive, got err: %v", r.Err)
	}
	if r.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestProbeClosedPort(t *testing.T) {
	c := New(200 * time.Millisecond)
	r := c.Probe("127.0.0.1", 1, "tcp")
	if r.Alive {
		t.Error("expected not alive")
	}
	if r.Err == nil {
		t.Error("expected error")
	}
}

func TestProbeAllReturnsAllResults(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	c := New(time.Second)
	targets := []Target{
		{Port: port, Proto: "tcp"},
		{Port: 1, Proto: "tcp"},
	}
	results := c.ProbeAll("127.0.0.1", targets)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !results[0].Alive {
		t.Error("first result should be alive")
	}
	if results[1].Alive {
		t.Error("second result should not be alive")
	}
}

func TestNewDefaultTimeout(t *testing.T) {
	c := New(0)
	if c.timeout != 2*time.Second {
		t.Errorf("expected 2s default, got %v", c.timeout)
	}
}
