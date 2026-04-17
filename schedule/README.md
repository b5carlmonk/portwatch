# schedule

The `schedule` package provides a simple interval-based task runner for portwatch.

## Overview

A `Runner` executes a user-supplied task function on a fixed interval. It runs the task **immediately** on start, then repeats at each tick until the context is cancelled.

## Usage

```go
runner := schedule.New(
    30*time.Second,
    func(ctx context.Context) error {
        // perform port scan, diff, alert
        return nil
    },
    func(err error) {
        log.Println("scan error:", err)
    },
)

ctx, cancel := context.WithCancel(context.Background())
defer cancel()
runner.Start(ctx) // blocks
```

## API

### `New(interval, task, onError) *Runner`

Creates a new Runner.

- `interval` — how often to run the task
- `task` — function called each tick; receives the context
- `onError` — called when task returns an error (optional, pass nil to ignore)

### `(*Runner).Start(ctx)`

Starts the runner. Blocks until `ctx` is cancelled.
