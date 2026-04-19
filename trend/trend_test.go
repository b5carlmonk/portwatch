package trend

import (
	"os"
	"path/filepath"
	"testing"
)

const testKey = "127.0.0.1:80/tcp"

func TestRecordOpenIncrementsCount(t *testing.T) {
	tr := New()
	tr.RecordOpen(testKey)
	tr.RecordOpen(testKey)
	e, ok := tr.Get(testKey)
	if !ok {
		t.Fatal("expected entry")
	}
	if e.OpenCount != 2 {
		t.Fatalf("expected 2 got %d", e.OpenCount)
	}
}

func TestRecordCloseIncrementsCount(t *testing.T) {
	tr := New()
	tr.RecordClose(testKey)
	e, _ := tr.Get(testKey)
	if e.CloseCount != 1 {
		t.Fatalf("expected 1 got %d", e.CloseCount)
	}
}

func TestGetMissingKey(t *testing.T) {
	tr := New()
	_, ok := tr.Get("missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestAllReturnsSnapshot(t *testing.T) {
	tr := New()
	tr.RecordOpen("a")
	tr.RecordOpen("b")
	if len(tr.All()) != 2 {
		t.Fatal("expected 2 entries")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tr := New()
	tr.RecordOpen(testKey)
	tr.RecordClose(testKey)

	dir := t.TempDir()
	p := filepath.Join(dir, "trend.json")
	if err := tr.Save(p); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	e, ok := loaded.Get(testKey)
	if !ok {
		t.Fatal("expected entry after load")
	}
	if e.OpenCount != 1 || e.CloseCount != 1 {
		t.Fatalf("unexpected counts: %+v", e)
	}
}

func TestLoadMissingFile(t *testing.T) {
	tr, err := Load(filepath.Join(t.TempDir(), "nope.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(tr.All()) != 0 {
		t.Fatal("expected empty tracker")
	}
}

func TestLastSeenUpdated(t *testing.T) {
	tr := New()
	tr.RecordOpen(testKey)
	e, _ := tr.Get(testKey)
	if e.LastSeen.IsZero() {
		t.Fatal("expected LastSeen to be set")
	}
}

func init() { _ = os.Getenv } // suppress unused import
