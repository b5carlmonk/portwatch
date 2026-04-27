# window

The `window` package provides a thread-safe sliding time-window counter for
tracking scan events and computing event rates over a configurable duration.

## Usage

```go
w := window.New(30 * time.Second)

// Record events.
w.Add(1)          // one new port change observed
w.Add(5)          // five more

// Query totals within the last 30 seconds.
fmt.Println(w.Total()) // sum of counts still inside the window
fmt.Println(w.Len())   // number of individual Add calls still inside the window

// Clear all accumulated data.
w.Reset()
```

## Behaviour

- Entries older than the configured duration are lazily evicted on every
  `Add`, `Total`, or `Len` call.
- All operations are safe for concurrent use.
- `Reset` immediately discards all entries regardless of age.

## Integration

Useful alongside `ratelimit`, `throttle`, and `cooldown` when you need a
running total rather than a simple allow/deny decision — for example, to
count how many port-change events occurred in the last minute before
deciding whether to escalate an alert.
