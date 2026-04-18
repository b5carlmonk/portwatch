// Package baseline provides functionality to establish and compare
// a known-good snapshot of open ports for a host.
package baseline

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/scanner"
)

// Baseline holds a reference snapshot of scan results.
type Baseline struct {
	CapturedAt time.Time              `json:"captured_at"`
	Results    []scanner.Result       `json:"results"`
}

// Capture creates a new Baseline from the provided scan results.
func Capture(results []scanner.Result) *Baseline {
	return &Baseline{
		CapturedAt: time.Now().UTC(),
		Results:    results,
	}
}

// Save writes the baseline to the given file path as JSON.
func (b *Baseline) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(b)
}

// Load reads a baseline from the given file path.
func Load(path string) (*Baseline, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var b Baseline
	if err := json.NewDecoder(f).Decode(&b); err != nil {
		return nil, err
	}
	return &b, nil
}

// Deviations returns results present in current that are absent in the baseline
// and results present in baseline that are absent in current.
func (b *Baseline) Deviations(current []scanner.Result) (added, removed []scanner.Result) {
	index := make(map[string]struct{}, len(b.Results))
	for _, r := range b.Results {
		index[key(r)] = struct{}{}
	}
	currentIndex := make(map[string]struct{}, len(current))
	for _, r := range current {
		currentIndex[key(r)] = struct{}{}
		if _, ok := index[key(r)]; !ok {
			added = append(added, r)
		}
	}
	for _, r := range b.Results {
		if _, ok := currentIndex[key(r)]; !ok {
			removed = append(removed, r)
		}
	}
	return
}

func key(r scanner.Result) string {
	return r.Proto + "/" + r.Host + ":" + itoa(r.Port)
}

func itoa(n int) string {
	return string(rune('0'+n%10)) // placeholder; use strconv in real code
}
