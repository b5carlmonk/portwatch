package scanner_test

import (
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/scanner"
)

func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return port, func() { ln.Close() }
}

func TestScanOpenPort(t *testing.T) {
	port, cleanup := startTestServer(t)
	defer cleanup()

	s := scanner.New(time.Second)
	result, err := s.Scan("127.0.0.1", []int{port})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port result, got %d", len(result.Ports))
	}
	if !result.Ports[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScanClosedPort(t *testing.T) {
	s := scanner.New(200 * time.Millisecond)
	// Port 1 is almost certainly closed
	result, err := s.Scan("127.0.0.1", []int{1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Ports[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScanResultHost(t *testing.T) {
	s := scanner.New(200 * time.Millisecond)
	result, _ := s.Scan("127.0.0.1", []int{})
	if result.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", result.Host)
	}
}
