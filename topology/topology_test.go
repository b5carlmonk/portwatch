package topology

import (
	"strings"
	"testing"

	"github.com/user/portwatch/scanner"
)

func makeResults(host string, ports ...int) []scanner.Result {
	results := make([]scanner.Result, 0, len(ports))
	for _, p := range ports {
		results = append(results, scanner.Result{
			Host:     host,
			Port:     p,
			Proto:    "tcp",
			Open:     true,
		})
	}
	return results
}

func TestUpdateStoresHost(t *testing.T) {
	m := New()
	m.Update(makeResults("192.168.1.1", 22, 80))

	n, ok := m.Get("192.168.1.1")
	if !ok {
		t.Fatal("expected host to be present")
	}
	if len(n.Ports) != 2 {
		t.Fatalf("expected 2 ports, got %d", len(n.Ports))
	}
}

func TestGetMissingHost(t *testing.T) {
	m := New()
	_, ok := m.Get("10.0.0.1")
	if ok {
		t.Fatal("expected missing host to return false")
	}
}

func TestHostsReturnsSorted(t *testing.T) {
	m := New()
	m.Update(makeResults("10.0.0.3", 80))
	m.Update(makeResults("10.0.0.1", 22))
	m.Update(makeResults("10.0.0.2", 443))

	hosts := m.Hosts()
	if len(hosts) != 3 {
		t.Fatalf("expected 3 hosts, got %d", len(hosts))
	}
	if hosts[0] != "10.0.0.1" || hosts[1] != "10.0.0.2" || hosts[2] != "10.0.0.3" {
		t.Fatalf("hosts not sorted: %v", hosts)
	}
}

func TestUpdateOverwritesPreviousPorts(t *testing.T) {
	m := New()
	m.Update(makeResults("192.168.1.1", 22, 80, 443))
	m.Update(makeResults("192.168.1.1", 8080))

	n, _ := m.Get("192.168.1.1")
	if len(n.Ports) != 1 {
		t.Fatalf("expected 1 port after update, got %d", len(n.Ports))
	}
	if n.Ports[0].Port != 8080 {
		t.Fatalf("expected port 8080, got %d", n.Ports[0].Port)
	}
}

func TestSummaryContainsHost(t *testing.T) {
	m := New()
	m.Update(makeResults("172.16.0.1", 22, 80))

	s := m.Summary()
	if !strings.Contains(s, "172.16.0.1") {
		t.Fatalf("summary missing host: %s", s)
	}
	if !strings.Contains(s, "2 port(s)") {
		t.Fatalf("summary missing port count: %s", s)
	}
}

func TestSummaryNoHosts(t *testing.T) {
	m := New()
	s := m.Summary()
	if !strings.Contains(s, "no hosts") {
		t.Fatalf("expected 'no hosts' in summary, got: %s", s)
	}
}
