// Package watchdog ties together scanning, diffing, filtering, alerting,
// and history into a single watch cycle.
package watchdog

import (
	"context"
	"time"

	"portwatch/alert"
	"portwatch/config"
	"portwatch/filter"
	"portwatch/history"
	"portwatch/report"
	"portwatch/scanner"
	"portwatch/state"
)

// Watchdog runs a single scan cycle: scan → filter → diff → alert → persist.
type Watchdog struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	filter  *filter.Filter
	alerter *alert.Alerter
	report  *report.Reporter
	hist    *history.History
}

// New creates a Watchdog from the provided config.
func New(cfg *config.Config, sc *scanner.Scanner, f *filter.Filter, a *alert.Alerter, r *report.Reporter, h *history.History) *Watchdog {
	return &Watchdog{cfg: cfg, scanner: sc, filter: f, alerter: a, report: r, hist: h}
}

// Run executes one watch cycle. It is safe to call from a scheduler tick.
func (w *Watchdog) Run(ctx context.Context) error {
	results, err := w.scanner.Scan(ctx, w.cfg.Host, w.cfg.Ports)
	if err != nil {
		return err
	}

	filtered := w.filter.Apply(results)

	prev, _ := state.Load(w.cfg.StateFile)
	diff := scanner.Diff(prev, filtered)

	if err := w.alerter.Notify(diff); err != nil {
		return err
	}

	w.report.PrintDiff(diff)

	w.hist.Add(history.Entry{
		Time:    time.Now(),
		Results: filtered,
		Diff:    diff,
	})
	if err := w.hist.Save(w.cfg.HistoryFile); err != nil {
		return err
	}

	return state.Save(w.cfg.StateFile, filtered)
}
