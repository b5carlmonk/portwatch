package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/scanner"
	"github.com/user/portwatch/snapshot"
)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "127.0.0.1", Port: 80, Protocol: "tcp", Open: true},
		{Host: "127.0.0.1", Port: 443, Protocol: "tcp", Open: true},
	}
}

func tmpStore(t *testing.T) *snapshot.Store {
	t.Helper()
	p := filepath.Join(t.TempDir(), "snapshots.json")
	return snapshot.New(p)
}

func TestAddAndGet(t *testing.T) {
	s := tmpStore(t)
	if err := s.Add("baseline", makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.Get("baseline")
	if !ok {
		t.Fatal("expected snapshot to exist")
	}
	if len(e.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(e.Results))
	}
}

func TestEmptyLabelErrors(t *testing.T) {
	s := tmpStore(t)
	if err := s.Add("", makeResults()); err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestDelete(t *testing.T) {
	s := tmpStore(t)
	_ = s.Add("tmp", makeResults())
	_ = s.Delete("tmp")
	if _, ok := s.Get("tmp"); ok {
		t.Fatal("expected snapshot to be deleted")
	}
}

func TestLabels(t *testing.T) {
	s := tmpStore(t)
	_ = s.Add("a", makeResults())
	_ = s.Add("b", makeResults())
	if len(s.Labels()) != 2 {
		t.Errorf("expected 2 labels, got %d", len(s.Labels()))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "snap.json")
	s1 := snapshot.New(p)
	_ = s1.Add("v1", makeResults())

	s2 := snapshot.New(p)
	if err := s2.Load(); err != nil {
		t.Fatalf("load error: %v", err)
	}
	if _, ok := s2.Get("v1"); !ok {
		t.Fatal("expected snapshot after reload")
	}
}

func TestLoadMissingFile(t *testing.T) {
	p := filepath.Join(t.TempDir(), "missing.json")
	s := snapshot.New(p)
	if err := s.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	_ = os.Remove(p)
}
