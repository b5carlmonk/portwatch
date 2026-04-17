# state

The `state` package provides persistence for port scan results between runs of portwatch.

## Overview

Each time portwatch runs a scan, the results can be saved as a **snapshot** — a timestamped JSON file on disk. On the next run, the previous snapshot is loaded and compared against the new scan using `scanner.Diff`, enabling detection of opened or closed ports since the last check.

## Usage

```go
// Save current scan results
err := state.Save("/var/lib/portwatch/state.json", results)

// Load previous snapshot
snap, err := state.Load("/var/lib/portwatch/state.json")
if err != nil {
    log.Fatal(err)
}

// If snap.Results is empty, this is the first run
diffs := scanner.Diff(snap.Results, currentResults)
```

## File Format

Snapshots are stored as pretty-printed JSON:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "results": [
    {"host": "localhost", "port": 80, "open": true, "service": "http"}
  ]
}
```
