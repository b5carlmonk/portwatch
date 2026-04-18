# digest

The `digest` package accumulates port scan snapshots over time and produces periodic summary reports showing which ports were observed and how frequently.

## Usage

```go
d := digest.New()           // writes to stdout
// or
d := digest.NewWithWriter(w) // writes to w

// after each scan cycle:
d.Add(results)

// periodically (e.g. daily):
d.Flush() // prints summary and clears state
```

## Output example

```
[digest] summary over 12 snapshots
  port 80/tcp      seen in 12 snapshot(s)
  port 443/tcp     seen in 12 snapshot(s)
  port 8080/tcp    seen in 3 snapshot(s)
```

## Integration

Wire `digest.Add` into the watchdog scan cycle and call `Flush` on a timer (e.g. via `schedule`) to emit daily or hourly summaries.
