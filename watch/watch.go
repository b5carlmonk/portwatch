// Package watch provides a high-level watcher that ties together scanning,
// diffing, filtering, and alerting into a single reusable component.
package watch

import (
	"context"
	"fmt"

	"github.com/user/portwatch/alert"
	"github.com/user/portwatch/scanner"
	"github.com/user/portwatch/state"
)

// Config holds the configuration for a Watcher.
type Config struct {
	// Targets is the list of host:port strings to scan.
	Targets []string
	// StateFile is the path where previous scan results are persisted.
	StateFile string
	// Timeout is the per-port dial timeout in milliseconds.
	TimeoutMs int
}

// Watcher runs a single watch cycle: scan → diff → alert.
type Watcher struct {
	cfg     Config
	scanner *scanner.Scanner
	alerter *alert.Alerter
}

// New creates a Watcher from cfg.
func New(cfg Config, a *alert.Alerter) (*Watcher, error) {
	if len(cfg.Targets) == 0 {
		return nil, fmt.Errorf("watch: at least one target is required")
	}
	s := scanner.New(cfg.Targets, cfg.TimeoutMs)
	return &Watcher{cfg: cfg, scanner: s, alerter: a}, nil
}

// Run performs one scan cycle and returns any diff-driven alerts.
// It loads previous state, scans, diffs, alerts, then persists new state.
func (w *Watcher) Run(ctx context.Context) error {
	prev, _ := state.Load(w.cfg.StateFile) // ignore missing-file error

	current, err := w.scanner.Scan(ctx)
	if err != nil {
		return fmt.Errorf("watch: scan failed: %w", err)
	}

	diff := scanner.Diff(prev, current)

	if err := w.alerter.Notify(diff); err != nil {
		return fmt.Errorf("watch: alert failed: %w", err)
	}

	if err := state.Save(w.cfg.StateFile, current); err != nil {
		return fmt.Errorf("watch: state save failed: %w", err)
	}
	return nil
}
