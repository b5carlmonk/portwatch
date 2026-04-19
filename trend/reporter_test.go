package trend

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintContainsKey(t *testing.T) {
	tr := New()
	tr.RecordOpen("192.168.1.1:443/tcp")

	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Print(tr)

	if !strings.Contains(buf.String(), "192.168.1.1:443/tcp") {
		t.Fatalf("expected key in output, got:\n%s", buf.String())
	}
}

func TestPrintShowsCounts(t *testing.T) {
	tr := New()
	tr.RecordOpen("host:22/tcp")
	tr.RecordOpen("host:22/tcp")
	tr.RecordClose("host:22/tcp")

	var buf bytes.Buffer
	NewReporter(&buf).Print(tr)
	out := buf.String()

	if !strings.Contains(out, "opened=2") {
		t.Errorf("expected opened=2 in output")
	}
	if !strings.Contains(out, "closed=1") {
		t.Errorf("expected closed=1 in output")
	}
}

func TestPrintEmptyTracker(t *testing.T) {
	var buf bytes.Buffer
	NewReporter(&buf).Print(New())
	if !strings.Contains(buf.String(), "(no data)") {
		t.Error("expected (no data) for empty tracker")
	}
}

func TestPrintSortedOutput(t *testing.T) {
	tr := New()
	tr.RecordOpen("z:9/tcp")
	tr.RecordOpen("a:1/tcp")

	var buf bytes.Buffer
	NewReporter(&buf).Print(tr)
	out := buf.String()

	iA := strings.Index(out, "a:1")
	iZ := strings.Index(out, "z:9")
	if iA > iZ {
		t.Error("expected sorted output")
	}
}
