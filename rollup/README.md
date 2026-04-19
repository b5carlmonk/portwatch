# rollup

The `rollup` package aggregates multiple scan diffs into a summary over a sliding time window.

## Usage

```go
r := rollup.New(5 * time.Minute)

// Add diffs as they arrive from the watchdog cycle.
r.Add(diff)

// Query totals for the current window.
opened, closed := r.Summary()
fmt.Printf("Last 5 min: +%d opened, -%d closed\n", opened, closed)

// Inspect individual entries.
for _, e := range r.Entries() {
    fmt.Println(e.At, len(e.Diff.Opened), len(e.Diff.Closed))
}

// Clear all buffered entries.
r.Flush()
```

## Behaviour

- Entries older than the configured window are pruned automatically on every `Add`, `Summary`, and `Entries` call.
- `Flush` immediately removes all entries regardless of age.
- All methods are safe for concurrent use.
