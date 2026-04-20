# correlate

The `correlate` package groups related port-change events across multiple hosts into **incidents** based on a configurable time window.

## Overview

When portwatch detects changes on many hosts in quick succession the raw stream of individual diffs can be noisy. `correlate` buffers those events and, on each flush, partitions them into incidents where consecutive events fall within a sliding time window.

## Usage

```go
import "github.com/user/portwatch/correlate"

c := correlate.New(10*time.Second, func(inc correlate.Incident) {
    fmt.Printf("Incident %s — %d related events\n", inc.ID, len(inc.Events))
    for _, e := range inc.Events {
        fmt.Printf("  %s:%d/%s %s\n", e.Host, e.Port, e.Proto, e.State)
    }
})

// Record events as diffs arrive
c.Add("192.168.1.1", 80,  "tcp", "opened")
c.Add("192.168.1.2", 80,  "tcp", "opened")
c.Add("192.168.1.3", 443, "tcp", "opened")

// Flush at the end of each scan cycle
c.Flush()
```

## Types

| Type | Description |
|------|-------------|
| `Event` | A single port change (host, port, proto, state, timestamp) |
| `Incident` | A slice of related `Event` values with a generated ID |
| `Correlator` | Buffers events and groups them on `Flush` |

## Methods

- `New(window, onIncident)` — create a correlator with a time window and callback
- `Add(host, port, proto, state)` — record a new event
- `Flush()` — group buffered events into incidents and fire the callback
- `Len()` — number of currently buffered events
