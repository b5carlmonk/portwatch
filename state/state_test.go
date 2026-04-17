package state

import (
	"os"
	"testing"

	"github.com/user/portwatch/scanner"
)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "localhost", Port: 80, Open: true, Service: "http"},
		{Host: "localhost", Port: 443, Open: true, Service: "https"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-state-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	results := makeResults()
	if err := Save(tmp.Name(), results); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	snap, err := Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(snap.Results) != len(results) {
		t.Errorf("expected %d results, got %d", len(results), len(snap.Results))
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	for i, r := range snap.Results {
		if r.Port != results[i].Port {
			t.Errorf("result[%d] port: expected %d, got %d", i, results[i].Port, r.Port)
		}
	}
}

func TestLoadMissingFile(t *testing.T) {
	snap, err := Load("/tmp/portwatch-nonexistent-state.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if len(snap.Results) != 0 {
		t.Errorf("expected empty results for missing file")
	}
	if !snap.Timestamp.IsZero() {
		t.Error("expected zero timestamp for missing file")
	}
}

func TestSaveEmptyResults(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-state-empty-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	if err := Save(tmp.Name(), []scanner.Result{}); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	snap, err := Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(snap.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(snap.Results))
	}
}
