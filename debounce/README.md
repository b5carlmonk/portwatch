# debounce

The `debounce` package reduces alert noise caused by transient port flaps by
requiring a change to be observed for a configurable number of **consecutive
scan cycles** before it is considered confirmed and forwarded to downstream
handlers.

## How it works

Each call to `Evaluate` receives a `scanner.Diff` (opened / closed port lists)
produced by the current scan cycle. Internally the debouncer maintains a
counter per unique change key (`direction:host:port:proto`). When the counter
reaches the configured threshold the change is included in the returned
`scanner.Diff` and its counter is removed. If a change disappears from the
diff before reaching the threshold its counter is reset.

## Usage

```go
import "github.com/user/portwatch/debounce"

// Require a port change to appear in 3 consecutive scans before alerting.
db := debounce.New(3)

// Inside your scan loop:
confirmed := db.Evaluate(diff)
if len(confirmed.Opened) > 0 || len(confirmed.Closed) > 0 {
    alert.Notify(confirmed)
}
```

## API

| Symbol | Description |
|---|---|
| `New(threshold int) *Debouncer` | Create a new debouncer with the given threshold (minimum 1). |
| `(*Debouncer).Evaluate(diff scanner.Diff) scanner.Diff` | Feed a diff; returns only confirmed changes. |
| `(*Debouncer).PendingLen() int` | Number of changes awaiting confirmation. |
| `(*Debouncer).Flush()` | Discard all pending unconfirmed entries. |
