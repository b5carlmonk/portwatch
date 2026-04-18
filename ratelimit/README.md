# ratelimit

The `ratelimit` package provides a simple in-memory rate limiter used to
suppress repeated alerts for the same port event within a configurable
time window.

## Usage

```go
limiter := ratelimit.New(5 * time.Minute)

if limiter.Allow("tcp:22") {
    // send alert — first occurrence or interval has elapsed
}
```

## API

| Function | Description |
|---|---|
| `New(interval)` | Create a limiter with the given suppression interval |
| `Allow(key) bool` | Returns true if the event should be processed |
| `Reset(key)` | Clear the timestamp for a single key |
| `Flush()` | Clear all recorded timestamps |

## Keys

Keys are arbitrary strings. By convention portwatch uses `"protocol:port"`
(e.g. `"tcp:443"`) to deduplicate alerts per service.

## Thread Safety

All methods are safe for concurrent use.
