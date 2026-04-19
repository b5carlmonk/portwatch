// Package watchdog orchestrates a full scan cycle.
package watchdog

import (
	"context"
	"log"

	"github.com/user/portwatch/alert"
	"github.com/user/portwatch/config"
	"github.com/user/portwatch/scanner"
	"github.com/user/portwatch/state"
)

// Watchdog runs periodic scan cycles.
type Watchdog struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	alert   *alert.Alert
}

// New returns a configured Watchdog.
func New(cfg *config.Config, sc *scanner.Scanner, al *alert.Alert) *Watchdog {
	return &Watchdog{cfg: cfg, scanner: sc, alert: al}
}

// RunCycle performs one scan, diffs against previous state, fires alerts and
// persists the new state.
func (w *Watchdog) RunCycle(ctx context.Context) error {
	results, err := w.scanner.Scan(ctx, w.cfg.Host, w.cfg.Ports)
	if err != nil {
		return err
	}

	prev, _ := state.Load(w.cfg.StateFile)
	diff := scanner.Diff(prev, results)

	if err := w.alert.Notify(diff); err != nil {
		log.Printf("alert error: %v", err)
	}

	return state.Save(w.cfg.StateFile, results)
}
