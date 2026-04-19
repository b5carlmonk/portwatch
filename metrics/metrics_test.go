package metrics

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestRecordScanIncrements(t *testing.T) {
	tr := NewWithWriter(&bytes.Buffer{})
	tr.RecordScan(10, 2, nil)
	tr.RecordScan(8, 0, nil)
	s := tr.Snapshot()
	if s.Scans != 2 {
		t.Fatalf("expected 2 scans, got %d", s.Scans)
	}
	if s.Changes != 2 {
		t.Fatalf("expected 2 total changes, got %d", s.Changes)
	}
}

func TestRecordScanUpdatesOpenPorts(t *testing.T) {
	tr := NewWithWriter(&bytes.Buffer{})
	tr.RecordScan(5, 0, nil)
	tr.RecordScan(7, 0, nil)
	s := tr.Snapshot()
	if s.OpenPorts != 7 {
		t.Fatalf("expected open=7, got %d", s.OpenPorts)
	}
}

func TestRecordScanCountsErrors(t *testing.T) {
	tr := NewWithWriter(&bytes.Buffer{})
	tr.RecordScan(0, 0, errors.New("boom"))
	tr.RecordScan(0, 0, nil)
	s := tr.Snapshot()
	if s.Errors != 1 {
		t.Fatalf("expected 1 error, got %d", s.Errors)
	}
}

func TestPrintContainsFields(t *testing.T) {
	var buf bytes.Buffer
	tr := NewWithWriter(&buf)
	tr.RecordScan(3, 1, nil)
	tr.Print()
	out := buf.String()
	for _, want := range []string{"scans=", "open=", "changes=", "errors=", "last="} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q: %s", want, out)
		}
	}
}

func TestSnapshotIsIndependent(t *testing.T) {
	tr := NewWithWriter(&bytes.Buffer{})
	tr.RecordScan(4, 1, nil)
	s1 := tr.Snapshot()
	tr.RecordScan(6, 2, nil)
	s2 := tr.Snapshot()
	if s1.Scans == s2.Scans {
		t.Fatal("snapshots should differ after additional scan")
	}
}
