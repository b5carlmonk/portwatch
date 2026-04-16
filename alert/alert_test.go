package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/alert"
	"github.com/user/portwatch/scanner"
)

func makeDiff(opened, closed []scanner.Result) scanner.DiffResult {
	return scanner.DiffResult{Opened: opened, Closed: closed}
}

func TestNotifyOpenedPort(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := makeDiff(
		[]scanner.Result{{Host: "localhost", Port: 8080, Service: "http"}},
		nil,
	)
	n.Notify(diff)

	out := buf.String()
	if !strings.Contains(out, "OPENED") {
		t.Errorf("expected OPENED in output, got: %s", out)
	}
	if !strings.Contains(out, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", out)
	}
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected ALERT level in output, got: %s", out)
	}
}

func TestNotifyClosedPort(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := makeDiff(
		nil,
		[]scanner.Result{{Host: "localhost", Port: 22, Service: "ssh"}},
	)
	n.Notify(diff)

	out := buf.String()
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED in output, got: %s", out)
	}
	if !strings.Contains(out, "WARN") {
		t.Errorf("expected WARN level in output, got: %s", out)
	}
}

func TestNotifyNoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := makeDiff(nil, nil)
	n.Notify(diff)

	if buf.Len() != 0 {
		t.Errorf("expected no output for empty diff, got: %s", buf.String())
	}
}
