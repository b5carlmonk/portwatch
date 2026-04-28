package anomaly

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintNilAlertWritesNothing(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporterWithWriter(&buf)
	r.Print(nil)
	if buf.Len() != 0 {
		t.Errorf("expected empty output for nil alert, got %q", buf.String())
	}
}

func TestPrintContainsHost(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporterWithWriter(&buf)
	r.Print(&Alert{Host: "192.168.1.1", Observed: 5, Mean: 2.0, StdDev: 0.5, Message: "test"})
	if !strings.Contains(buf.String(), "192.168.1.1") {
		t.Errorf("expected host in output, got %q", buf.String())
	}
}

func TestPrintContainsObserved(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporterWithWriter(&buf)
	r.Print(&Alert{Host: "h", Observed: 42, Mean: 3.0, StdDev: 1.0, Message: "msg"})
	if !strings.Contains(buf.String(), "42") {
		t.Errorf("expected observed count in output, got %q", buf.String())
	}
}

func TestPrintContainsMessage(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporterWithWriter(&buf)
	a := &Alert{Host: "h", Observed: 10, Mean: 2.0, StdDev: 1.0, Message: "unique-msg-xyz"}
	r.Print(a)
	if !strings.Contains(buf.String(), "unique-msg-xyz") {
		t.Errorf("expected message in output, got %q", buf.String())
	}
}

func TestPrintContainsANOMALYLabel(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporterWithWriter(&buf)
	r.Print(&Alert{Host: "h", Observed: 1, Mean: 1.0, StdDev: 0.1, Message: "m"})
	if !strings.Contains(buf.String(), "ANOMALY") {
		t.Errorf("expected ANOMALY label in output, got %q", buf.String())
	}
}
