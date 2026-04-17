package report

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/scanner"
)

func makeResults(host string, ports []int, open bool) []scanner.Result {
	var results []scanner.Result
	for _, p := range ports {
		results = append(results, scanner.Result{Host: host, Port: p, Proto: "tcp", Open: open})
	}
	return results
}

func TestPrintResultsContainsPorts(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	results := makeResults("localhost", []int{80, 443}, true)
	r.PrintResults(results, time.Now())
	out := buf.String()
	if !strings.Contains(out, "80") || !strings.Contains(out, "443") {
		t.Errorf("expected ports in output, got: %s", out)
	}
}

func TestPrintResultsHost(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	results := makeResults("myhost", []int{22}, true)
	r.PrintResults(results, time.Now())
	if !strings.Contains(buf.String(), "myhost") {
		t.Errorf("expected host in output")
	}
}

func TestPrintDiffOpened(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	diff := scanner.Diff{
		Opened: makeResults("localhost", []int{8080}, true),
	}
	r.PrintDiff(diff)
	if !strings.Contains(buf.String(), "OPEN") {
		t.Errorf("expected OPEN in diff output")
	}
}

func TestPrintDiffClosed(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	diff := scanner.Diff{
		Closed: makeResults("localhost", []int{3306}, false),
	}
	r.PrintDiff(diff)
	if !strings.Contains(buf.String(), "CLOSED") {
		t.Errorf("expected CLOSED in diff output")
	}
}

func TestPrintDiffNoChanges(t *testing.T) {
	var buf bytes.Buffer
	r := New(&buf)
	r.PrintDiff(scanner.Diff{})
	if !strings.Contains(buf.String(), "No changes") {
		t.Errorf("expected no-changes message")
	}
}
