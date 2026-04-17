package schedule

import (
	"context"
	"time"
)

// Runner executes a function on a fixed interval until the context is cancelled.
type Runner struct {
	interval time.Duration
	task     func(ctx context.Context) error
	onError  func(err error)
}

// New creates a new Runner with the given interval and task.
func New(interval time.Duration, task func(ctx context.Context) error, onError func(err error)) *Runner {
	if onError == nil {
		onError = func(err error) {}
	}
	return &Runner{
		interval: interval,
		task:     task,
		onError:  onError,
	}
}

// Start begins executing the task at each interval tick.
// It runs the task immediately on start, then waits for each tick.
// Blocks until ctx is cancelled.
func (r *Runner) Start(ctx context.Context) {
	if err := r.task(ctx); err != nil {
		r.onError(err)
	}

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.task(ctx); err != nil {
				r.onError(err)
			}
		}
	}
}
