package history_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/history"
	"github.com/user/portwatch/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	var results []scanner.Result
	for _, p := range ports {
		results = append(results, scanner.Result{
			Host:     "localhost",
			Port:     p,
			Protocol: "tcp",
			Open:     true,
		})
	}
	return results
}

func TestAddEntry(t *testing.T) {
	h := history.New(10)
	h.Add(makeResults(80, 443))
	if len(h.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(h.Entries))
	}
	if len(h.Entries[0].Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(h.Entries[0].Results))
	}
}

func TestMaxSizeEviction(t *testing.T) {
	h := history.New(3)
	for i := 0; i < 5; i++ {
		h.Add(makeResults(8080))
	}
	if len(h.Entries) != 3 {
		t.Fatalf("expected 3 entries after eviction, got %d", len(h.Entries))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	h := history.New(10)
	h.Add(makeResults(22, 80))
	h.Add(makeResults(443))

	if err := h.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := history.Load(path, 10)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(loaded.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(loaded.Entries))
	}
	// Verify results are preserved across save/load.
	if len(loaded.Entries[0].Results) != 2 {
		t.Fatalf("expected 2 results in first entry, got %d", len(loaded.Entries[0].Results))
	}
	if len(loaded.Entries[1].Results) != 1 {
		t.Fatalf("expected 1 result in second entry, got %d", len(loaded.Entries[1].Results))
	}
}

func TestLoadMissingFile(t *testing.T) {
	h, err := history.Load("/nonexistent/path/history.json", 10)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(h.Entries) != 0 {
		t.Fatalf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestTimestampSet(t *testing.T) {
	h := history.New(10)
	h.Add(makeResults(80))
	if h.Entries[0].Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestDefaultMaxSize(t *testing.T) {
	h := history.New(0)
	for i := 0; i < 110; i++ {
		h.Add(makeResults(80))
	}
	if len(h.Entries) > 100 {
		t.Fatalf("expected at most 100 entries, got %d", len(h.Entries))
	}
	_ = os.Getenv("") // suppress unused import warning
}
