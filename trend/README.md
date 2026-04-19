# trend

The `trend` package tracks how often individual ports are opened or closed over time.

## Overview

Each port is identified by a string key in the form `host:port/proto`. The `Tracker` records open and close events and persists the data as JSON.

## Usage

```go
tr, _ := trend.Load("/var/lib/portwatch/trend.json")

// Record events from a diff
for _, d := range diffs {
    if d.State == "opened" {
        tr.RecordOpen(d.Key)
    } else {
        tr.RecordClose(d.Key)
    }
}

// Print summary
trend.NewReporter(os.Stdout).Print(tr)

// Persist
tr.Save("/var/lib/portwatch/trend.json")
```

## Types

### Tracker

| Method | Description |
|---|---|
| `RecordOpen(key)` | Increment open counter |
| `RecordClose(key)` | Increment close counter |
| `Get(key)` | Retrieve entry by key |
| `All()` | Return snapshot of all entries |
| `Save(path)` | Persist to JSON file |
| `Load(path)` | Restore from JSON file |

### Reporter

Prints a sorted human-readable summary of all tracked ports.
