package shadow_test

import (
	"testing"

	"github.com/user/portwatch/scanner"
	"github.com/user/portwatch/shadow"
)

func makeResults(host string, ports []int, open bool) []scanner.Result {
	var out []scanner.Result
	for _, p := range ports {
		out = append(out, scanner.Result{
			Host:     host,
			Port:     p,
			Protocol: "tcp",
			Open:     open,
		})
	}
	return out
}

func TestNoShadowWhenAllInBaseline(t *testing.T) {
	baseline := makeResults("localhost", []int{22, 80, 443}, true)
	d := shadow.New(baseline)
	findings := d.Detect(makeResults("localhost", []int{22, 80, 443}, true))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestDetectsShadowPort(t *testing.T) {
	baseline := makeResults("localhost", []int{22, 80}, true)
	d := shadow.New(baseline)
	current := makeResults("localhost", []int{22, 80, 9999}, true)
	findings := d.Detect(current)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Port != 9999 {
		t.Errorf("expected port 9999, got %d", findings[0].Port)
	}
}

func TestClosedPortsIgnored(t *testing.T) {
	baseline := makeResults("localhost", []int{22}, true)
	d := shadow.New(baseline)
	// Port 9999 is not in baseline but is closed — should not be flagged.
	current := append(
		makeResults("localhost", []int{22}, true),
		makeResults("localhost", []int{9999}, false)...,
	)
	findings := d.Detect(current)
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestAllowAddsToBaseline(t *testing.T) {
	baseline := makeResults("localhost", []int{22}, true)
	d := shadow.New(baseline)
	d.Allow(scanner.Result{Host: "localhost", Port: 8080, Protocol: "tcp", Open: true})
	if d.Len() != 2 {
		t.Fatalf("expected len 2, got %d", d.Len())
	}
	findings := d.Detect(makeResults("localhost", []int{22, 8080}, true))
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings after Allow, got %d", len(findings))
	}
}

func TestMultipleShadowPorts(t *testing.T) {
	baseline := makeResults("host1", []int{80}, true)
	d := shadow.New(baseline)
	current := makeResults("host1", []int{80, 8080, 8443}, true)
	findings := d.Detect(current)
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
}

func TestBaselineClosedPortsNotAllowed(t *testing.T) {
	// Closed ports in the baseline should not be treated as allowed.
	baseline := makeResults("localhost", []int{22}, false)
	d := shadow.New(baseline)
	if d.Len() != 0 {
		t.Fatalf("expected empty allowed set, got %d", d.Len())
	}
	findings := d.Detect(makeResults("localhost", []int{22}, true))
	if len(findings) != 1 {
		t.Fatalf("expected 1 shadow finding, got %d", len(findings))
	}
}
