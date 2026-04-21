// Package watch orchestrates a complete port-watch cycle.
//
// A Watcher loads the previous scan state, performs a fresh scan of the
// configured targets, computes the diff, fires alerts for any changes, and
// persists the new state for the next cycle.
//
// Typical usage:
//
//	w, err := watch.New(cfg, alerter)
//	if err != nil { ... }
//	if err := w.Run(ctx); err != nil { ... }
package watch
