package baseline_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/baseline"
	"github.com/user/portwatch/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	results := make([]scanner.Result, len(ports))
	for i, p := range ports {
		results[i] = scanner.Result{Host: "localhost", Port: p, Proto: "tcp", Open: true}
	}
	return results
}

func TestCaptureStoresResults(t *testing.T) {
	res := makeResults(80, 443)
	b := baseline.Capture(res)
	if len(b.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(b.Results))
	}
	if b.CapturedAt.IsZero() {
		t.Fatal("expected non-zero CapturedAt")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	original := baseline.Capture(makeResults(22, 80))
	if err := original.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := baseline.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected non-nil baseline")
	}
	if len(loaded.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(loaded.Results))
	}
}

func TestLoadMissingFile(t *testing.T) {
	b, err := baseline.Load("/nonexistent/baseline.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b != nil {
		t.Fatal("expected nil for missing file")
	}
}

func TestDeviationsAdded(t *testing.T) {
	b := baseline.Capture(makeResults(80))
	current := makeResults(80, 8080)
	added, removed := b.Deviations(current)
	if len(added) != 1 {
		t.Fatalf("expected 1 added, got %d", len(added))
	}
	if len(removed) != 0 {
		t.Fatalf("expected 0 removed, got %d", len(removed))
	}
}

func TestDeviationsRemoved(t *testing.T) {
	b := baseline.Capture(makeResults(80, 443))
	current := makeResults(80)
	added, removed := b.Deviations(current)
	if len(added) != 0 {
		t.Fatalf("expected 0 added, got %d", len(added))
	}
	if len(removed) != 1 {
		t.Fatalf("expected 1 removed, got %d", len(removed))
	}
}

func TestDeviationsNoChange(t *testing.T) {
	b := baseline.Capture(makeResults(80, 443))
	current := makeResults(80, 443)
	added, removed := b.Deviations(current)
	if len(added) != 0 || len(removed) != 0 {
		t.Fatal("expected no deviations")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
