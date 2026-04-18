# history

The `history` package records scan results over time, enabling trend analysis and audit trails for port changes on a host.

## Usage

```go
h := history.New(100) // retain last 100 scans

// Add results from a scan
h.Add(results)

// Persist to disk
if err := h.Save("/var/lib/portwatch/history.json"); err != nil {
    log.Fatal(err)
}

// Load from disk
h, err := history.Load("/var/lib/portwatch/history.json", 100)
if err != nil {
    log.Fatal(err)
}
```

## Entry

Each `Entry` contains:

| Field       | Type              | Description                        |
|-------------|-------------------|------------------------------------|
| `Timestamp` | `time.Time`       | UTC time the scan was recorded     |
| `Results`   | `[]scanner.Result`| Port scan results at that instant  |

## Retention

When the number of entries exceeds `maxSize`, the oldest entries are evicted automatically. A `maxSize` of `0` or less defaults to `100`.
