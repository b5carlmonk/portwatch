# audit

The `audit` package provides a persistent, timestamped log of every port-scan cycle.

## Usage

```go
log := audit.New("/var/lib/portwatch/audit.json")

// Load previous entries from disk (optional)
_ = log.Load()

// After each scan cycle
log.Record("myhost", results)
_ = log.Save()

// Inspect entries
for _, e := range log.Entries() {
    fmt.Println(e.Time, e.Host, len(e.Results))
}
```

## Entry fields

| Field     | Type              | Description                        |
|-----------|-------------------|------------------------------------|
| `Time`    | `time.Time`       | UTC timestamp of the scan cycle    |
| `Host`    | `string`          | Target host that was scanned       |
| `Results` | `[]scanner.Result`| Full list of results for the cycle |

## Persistence

Entries are stored as a JSON array. Call `Save()` after each cycle and `Load()` on startup to retain history across restarts.
