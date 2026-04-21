# quota

The `quota` package enforces per-key scan rate limits using a sliding time window.

## Overview

A `Quota` allows at most `max` events for a given key within a rolling `window` duration. Once the limit is reached, further calls to `Allow` return `false` and a descriptive error until old events fall outside the window.

## Usage

```go
q := quota.New(5, time.Minute)

ok, err := q.Allow("192.168.1.1")
if !ok {
    log.Println("scan rate limit reached:", err)
    return
}
// proceed with scan
```

## API

| Function | Description |
|---|---|
| `New(max int, window time.Duration)` | Create a new Quota enforcer |
| `Allow(key string) (bool, error)` | Check and record an event for key |
| `Remaining(key string) int` | How many events remain in the current window |
| `Reset(key string)` | Clear quota state for a specific key |
| `Flush()` | Clear all quota state |

## Notes

- All methods are safe for concurrent use.
- Timestamps outside the window are evicted lazily on each `Allow` call.
- Keys are arbitrary strings — typically hostnames or IP addresses.
