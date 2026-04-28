package anomaly

import (
	"testing"

	"github.com/user/portwatch/scanner"
)

func makeResults(host string, ports []int) []scanner.Result {
	results := make([]scanner.Result, len(ports))
	for i, p := range ports {
		results[i] = scanner.Result{Host: host, Port: p, Proto: "tcp", Open: true}
	}
	return results
}

func TestInsufficientHistoryReturnsNil(t *testing.T) {
	d := New(2.0)
	r := makeResults("localhost", []int{80, 443})
	if got := d.Analyze(r); got != nil {
		t.Fatalf("expected nil with only 1 sample, got %+v", got)
	}
	if got := d.Analyze(r); got != nil {
		t.Fatalf("expected nil with only 2 samples, got %+v", got)
	}
}

func TestNoAnomalyWhenStable(t *testing.T) {
	d := New(2.0)
	r := makeResults("localhost", []int{80, 443})
	for i := 0; i < 10; i++ {
		if alert := d.Analyze(r); alert != nil {
			t.Fatalf("unexpected alert on stable data: %s", alert.Message)
		}
	}
}

func TestDetectsSpike(t *testing.T) {
	d := New(2.0)
	normal := makeResults("10.0.0.1", []int{80, 443})
	// Seed stable baseline.
	for i := 0; i < 10; i++ {
		d.Analyze(normal)
	}
	// Inject a spike: 20 open ports.
	spike := makeResults("10.0.0.1", make([]int, 20))
	for i := range spike {
		spike[i].Port = 1000 + i
		spike[i].Open = true
	}
	alert := d.Analyze(spike)
	if alert == nil {
		t.Fatal("expected anomaly alert for spike, got nil")
	}
	if alert.Observed != 20 {
		t.Errorf("expected Observed=20, got %d", alert.Observed)
	}
	if alert.Host != "10.0.0.1" {
		t.Errorf("expected Host=10.0.0.1, got %s", alert.Host)
	}
}

func TestResetClearsSamples(t *testing.T) {
	d := New(2.0)
	r := makeResults("h", []int{22})
	d.Analyze(r)
	d.Analyze(r)
	d.Analyze(r)
	d.Reset()
	// After reset, need 3 samples again before analysis fires.
	if got := d.Analyze(r); got != nil {
		t.Fatalf("expected nil after reset, got %+v", got)
	}
}

func TestZeroThresholdDefaultsToTwo(t *testing.T) {
	d := New(0)
	if d.threshold != 2.0 {
		t.Errorf("expected default threshold 2.0, got %f", d.threshold)
	}
}

func TestAlertMessageNonEmpty(t *testing.T) {
	d := New(2.0)
	normal := makeResults("h", []int{80, 443})
	for i := 0; i < 10; i++ {
		d.Analyze(normal)
	}
	spike := makeResults("h", make([]int, 25))
	for i := range spike {
		spike[i].Port = 2000 + i
		spike[i].Open = true
	}
	alert := d.Analyze(spike)
	if alert == nil {
		t.Fatal("expected alert")
	}
	if alert.Message == "" {
		t.Error("expected non-empty alert message")
	}
}
