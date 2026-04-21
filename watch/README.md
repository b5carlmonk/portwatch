# watch

The `watch` package provides a high-level orchestration layer that ties together
scanning, diffing, and alerting into a single reusable `Watcher` component.

## Overview

A single call to `Watcher.Run` performs the full cycle:

1. **Load** the previous scan state from disk (if any).
2. **Scan** all configured targets.
3. **Diff** the current results against the previous state.
4. **Alert** on any opened or closed ports.
5. **Persist** the new state for the next run.

## Usage

```go
import (
    "context"
    "os"

    "github.com/user/portwatch/alert"
    "github.com/user/portwatch/watch"
)

a := alert.New(os.Stdout)

cfg := watch.Config{
    Targets:   []string{"192.168.1.1:22", "192.168.1.1:80"},
    StateFile: "/var/lib/portwatch/state.json",
    TimeoutMs: 500,
}

w, err := watch.New(cfg, a)
if err != nil {
    log.Fatal(err)
}

if err := w.Run(context.Background()); err != nil {
    log.Fatal(err)
}
```

## Config

| Field | Description |
|-------|-------------|
| `Targets` | List of `host:port` strings to scan (required). |
| `StateFile` | Path to the JSON file used to persist scan state. |
| `TimeoutMs` | Per-port dial timeout in milliseconds. |
