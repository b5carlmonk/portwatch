package audit_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/audit"
	"github.com/user/portwatch/scanner"
)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "localhost", Port: 80, Proto: "tcp", Open: true},
		{Host: "localhost", Port: 443, Proto: "tcp", Open: true},
	}
}

func TestRecordAddsEntry(t *testing.T) {
	l := audit.New("")
	l.Record("localhost", makeResults())
	if len(l.Entries()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.Entries()))
	}
}

func TestRecordSetsHost(t *testing.T) {
	l := audit.New("")
	l.Record("myhost", makeResults())
	if l.Entries()[0].Host != "myhost" {
		t.Errorf("unexpected host: %s", l.Entries()[0].Host)
	}
}

func TestRecordSetsResults(t *testing.T) {
	l := audit.New("")
	res := makeResults()
	l.Record("localhost", res)
	if len(l.Entries()[0].Results) != len(res) {
		t.Errorf("results length mismatch")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	l := audit.New(path)
	l.Record("localhost", makeResults())
	l.Record("localhost", makeResults())
	if err := l.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}

	l2 := audit.New(path)
	if err := l2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(l2.Entries()) != 2 {
		t.Errorf("expected 2 entries after load, got %d", len(l2.Entries()))
	}
}

func TestLoadMissingFile(t *testing.T) {
	l := audit.New("/nonexistent/path/audit.json")
	if err := l.Load(); err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestEntriesTimestampSet(t *testing.T) {
	l := audit.New("")
	l.Record("h", makeResults())
	if l.Entries()[0].Time.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestSaveError(t *testing.T) {
	l := audit.New(filepath.Join(t.TempDir(), "sub", "audit.json"))
	l.Record("h", makeResults())
	if err := l.Save(); err == nil {
		t.Error("expected error saving to bad path")
	}
	_ = os.Remove("")
}
