package priority_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/priority"
)

func TestPrintContainsHeader(t *testing.T) {
	var buf strings.Builder
	r := priority.NewReporterWithWriter(&buf)
	r.Print(map[string]priority.Level{"22/tcp": priority.Critical})
	if !strings.Contains(buf.String(), "PORT/PROTO") {
		t.Error("expected header PORT/PROTO in output")
	}
}

func TestPrintContainsEntry(t *testing.T) {
	var buf strings.Builder
	r := priority.NewReporterWithWriter(&buf)
	r.Print(map[string]priority.Level{"443/tcp": priority.High})
	if !strings.Contains(buf.String(), "443/tcp") {
		t.Error("expected 443/tcp in output")
	}
	if !strings.Contains(buf.String(), "HIGH") {
		t.Error("expected HIGH in output")
	}
}

func TestPrintEmptyWritesNothing(t *testing.T) {
	var buf strings.Builder
	r := priority.NewReporterWithWriter(&buf)
	r.Print(nil)
	if buf.Len() != 0 {
		t.Errorf("expected empty output for nil scores, got %q", buf.String())
	}
}

func TestPrintSortedOutput(t *testing.T) {
	var buf strings.Builder
	r := priority.NewReporterWithWriter(&buf)
	r.Print(map[string]priority.Level{
		"9000/tcp": priority.Low,
		"22/tcp":   priority.Critical,
		"80/tcp":   priority.Medium,
	})
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + separator + 3 entries = 5 lines
	if len(lines) != 5 {
		t.Fatalf("expected 5 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[2], "22/tcp") {
		t.Errorf("expected 22/tcp first in sorted output, got %q", lines[2])
	}
}
