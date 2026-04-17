package state

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/scanner"
)

// Snapshot holds a saved scan result with metadata.
type Snapshot struct {
	Timestamp time.Time        `json:"timestamp"`
	Results   []scanner.Result `json:"results"`
}

// Save writes a snapshot of scan results to the given file path.
func Save(path string, results []scanner.Result) error {
	snap := Snapshot{
		Timestamp: time.Now(),
		Results:   results,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Load reads a snapshot from the given file path.
// Returns an empty Snapshot and no error if the file does not exist.
func Load(path string) (Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Snapshot{}, nil
		}
		return Snapshot{}, err
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return Snapshot{}, err
	}
	return snap, nil
}
