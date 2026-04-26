package scorecard_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/scorecard"
)

func TestPrintContainsHost(t *testing.T) {
	var buf bytes.Buffer
	r := scorecard.NewReporterWithWriter(&buf)
	rep := scorecard.Report{Host: "10.0.0.5", Score: 10, Level: scorecard.Low, OpenPorts: 1}
	r.Print(rep)
	if !strings.Contains(buf.String(), "10.0.0.5") {
		t.Error("expected host in output")
	}
}

func TestPrintContainsScore(t *testing.T) {
	var buf bytes.Buffer
	r := scorecard.NewReporterWithWriter(&buf)
	rep := scorecard.Report{Host: "h", Score: 42, Level: scorecard.Medium, OpenPorts: 3}
	r.Print(rep)
	if !strings.Contains(buf.String(), "42") {
		t.Error("expected score 42 in output")
	}
}

func TestPrintContainsLevel(t *testing.T) {
	var buf bytes.Buffer
	r := scorecard.NewReporterWithWriter(&buf)
	rep := scorecard.Report{Host: "h", Score: 90, Level: scorecard.Critical, OpenPorts: 5}
	r.Print(rep)
	if !strings.Contains(buf.String(), "CRITICAL") {
		t.Error("expected CRITICAL in output")
	}
}

func TestPrintBreakdownShown(t *testing.T) {
	var buf bytes.Buffer
	r := scorecard.NewReporterWithWriter(&buf)
	rep := scorecard.Report{
		Host:      "h",
		Score:     40,
		Level:     scorecard.High,
		OpenPorts: 1,
		Breakdown: []string{"port 23 (tcp) +40"},
	}
	r.Print(rep)
	if !strings.Contains(buf.String(), "port 23") {
		t.Error("expected breakdown entry in output")
	}
}
