# retry

The `retry` package provides a configurable retry mechanism with optional
exponential backoff for operations that may transiently fail.

## Usage

```go
import "portwatch/retry"

p := retry.Default() // 3 attempts, 500 ms base delay, ×2 backoff

err := p.Do(ctx, func() error {
    return doSomethingFallible()
})
if errors.Is(err, retry.ErrExhausted) {
    log.Println("all attempts failed")
}
```

## Policy fields

| Field | Description | Default |
|---|---|---|
| `MaxAttempts` | Total attempts including the first | `3` |
| `Delay` | Base wait between attempts | `500ms` |
| `Multiplier` | Delay growth factor per attempt | `2.0` |

Set `Multiplier` to `1.0` for a constant delay.

## Context support

`Do` checks `ctx` before each attempt and between waits, so cancellation
or deadline expiry is honoured promptly.

## Errors

- `ErrExhausted` — all attempts consumed; wraps the last underlying error.
- Context errors are returned directly without wrapping.
