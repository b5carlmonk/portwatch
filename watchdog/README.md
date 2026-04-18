# watchdog

The `watchdog` package orchestrates a complete port-watch cycle by composing
the lower-level `scanner`, `filter`, `alert`, `report`, `history`, and `state`
packages.

## Usage

```go
wd := watchdog.New(cfg, sc, f, alerter, reporter, hist)

// run once
if err := wd.Run(ctx); err != nil {
    log.Println("cycle error:", err)
}

// or drive it from a scheduler
runner := schedule.New(cfg.Interval, wd.Run, onError)
runner.Start(ctx)
```

## Cycle steps

| Step | Package |
|------|---------|
| Scan ports | `scanner` |
| Apply rules | `filter` |
| Diff vs previous | `scanner.Diff` |
| Send alerts | `alert` |
| Print report | `report` |
| Append history | `history` |
| Save state | `state` |

## Notes

- A missing state file on the first run is treated as an empty baseline — no
  false-positive "closed" alerts are generated.
- History eviction is controlled by the `MaxSize` set on the `history.History`
  instance passed to `New`.
