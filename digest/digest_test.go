package digest_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/digest"
	"github.com/user/portwatch/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	var ps []scanner.Port
	for _, p := range ports {
		ps = append(ps, scanner.Port{Port: p, Proto: "tcp"})
	}
	return []scanner.Result{{Host: "127.0.0.1", Ports: ps}}
}

func TestFlushEmptyWritesNoData(t *testing.T) {
	var buf bytes.Buffer
	d := digest.NewWithWriter(&buf)
	d.Flush()
	if !strings.Contains(buf.String(), "no data") {
		t.Errorf("expected no-data message, got: %s", buf.String())
	}
}

func TestAddIncreasesLen(t *testing.T) {
	d := digest.New()
	d.Add(makeResults(80))
	d.Add(makeResults(443))
	if d.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", d.Len())
	}
}

func TestFlushWritesSummary(t *testing.T) {
	var buf bytes.Buffer
	d := digest.NewWithWriter(&buf)
	d.Add(makeResults(80, 443))
	d.Add(makeResults(80))
	d.Flush()
	out := buf.String()
	if !strings.Contains(out, "2 snapshots") {
		t.Errorf("expected snapshot count, got: %s", out)
	}
	if !strings.Contains(out, "80/tcp") {
		t.Errorf("expected port 80/tcp in output, got: %s", out)
	}
	if !strings.Contains(out, "443/tcp") {
		t.Errorf("expected port 443/tcp in output, got: %s", out)
	}
}

func TestFlushClearsEntries(t *testing.T) {
	var buf bytes.Buffer
	d := digest.NewWithWriter(&buf)
	d.Add(makeResults(22))
	d.Flush()
	if d.Len() != 0 {
		t.Errorf("expected 0 entries after flush, got %d", d.Len())
	}
}

func TestPortSeenCount(t *testing.T) {
	var buf bytes.Buffer
	d := digest.NewWithWriter(&buf)
	d.Add(makeResults(8080))
	d.Add(makeResults(8080))
	d.Add(makeResults(8080))
	d.Flush()
	if !strings.Contains(buf.String(), "3 snapshot(s)") {
		t.Errorf("expected port seen 3 times, got: %s", buf.String())
	}
}
