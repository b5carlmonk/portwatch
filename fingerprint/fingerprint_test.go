package fingerprint_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/fingerprint"
	"github.com/user/portwatch/scanner"
)

// startBannerServer opens a TCP listener that writes banner on connect.
func startBannerServer(t *testing.T, banner string) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	go func() {
		defer ln.Close()
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Fprint(conn, banner)
	}()
	return port
}

func makeResult(host string, port int, open bool) scanner.Result {
	return scanner.Result{Host: host, Port: port, Proto: "tcp", Open: open}
}

func TestIdentifyRecognisesSSH(t *testing.T) {
	port := startBannerServer(t, "SSH-2.0-OpenSSH_8.9")
	f := fingerprint.New(500 * time.Millisecond)
	results := []scanner.Result{makeResult("127.0.0.1", port, true)}
	svcs := f.Identify(results)
	if len(svcs) != 1 {
		t.Fatalf("expected 1 service, got %d", len(svcs))
	}
	if svcs[0].Service != "ssh" {
		t.Errorf("expected ssh, got %q", svcs[0].Service)
	}
}

func TestIdentifySkipsClosedPorts(t *testing.T) {
	f := fingerprint.New(200 * time.Millisecond)
	results := []scanner.Result{makeResult("127.0.0.1", 9, false)}
	svcs := f.Identify(results)
	if len(svcs) != 0 {
		t.Errorf("expected 0 services, got %d", len(svcs))
	}
}

func TestIdentifyUnknownBanner(t *testing.T) {
	port := startBannerServer(t, "MYSTERY-PROTOCOL/1.0")
	f := fingerprint.New(500 * time.Millisecond)
	results := []scanner.Result{makeResult("127.0.0.1", port, true)}
	svcs := f.Identify(results)
	if len(svcs) != 1 {
		t.Fatalf("expected 1 service, got %d", len(svcs))
	}
	if svcs[0].Service != "unknown" {
		t.Errorf("expected unknown, got %q", svcs[0].Service)
	}
}

func TestIdentifyNoBannerEmptyService(t *testing.T) {
	// port that closes immediately without sending data
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	port := ln.Addr().(*net.TCPAddr).Port
	go func() { c, _ := ln.Accept(); c.Close(); ln.Close() }()

	f := fingerprint.New(300 * time.Millisecond)
	results := []scanner.Result{makeResult("127.0.0.1", port, true)}
	svcs := f.Identify(results)
	if len(svcs) != 1 {
		t.Fatalf("expected 1 service, got %d", len(svcs))
	}
	if svcs[0].Service != "" {
		t.Errorf("expected empty service, got %q", svcs[0].Service)
	}
}

func TestIdentifyPortFieldSet(t *testing.T) {
	port := startBannerServer(t, "SSH-2.0-test")
	f := fingerprint.New(500 * time.Millisecond)
	results := []scanner.Result{makeResult("127.0.0.1", port, true)}
	svcs := f.Identify(results)
	if svcs[0].Port != port {
		t.Errorf("expected port %d, got %d", port, svcs[0].Port)
	}
}
