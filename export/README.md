# export

The `export` package serialises scan results to common output formats for
downstream processing or archival.

## Supported formats

| Format | Description |
|--------|-------------|
| `json` | Pretty-printed JSON array of result objects |
| `csv`  | Comma-separated values with a header row |

## Usage

```go
import "github.com/user/portwatch/export"

// Write results as JSON to stdout
e := export.New(os.Stdout, export.FormatJSON)
if err := e.Write(results); err != nil {
    log.Fatal(err)
}

// Write results as CSV to a file
f, _ := os.Create("results.csv")
defer f.Close()
e = export.New(f, export.FormatCSV)
_ = e.Write(results)
```

## JSON schema

Each element in the JSON array contains:

```json
{
  "host": "localhost",
  "port": 443,
  "protocol": "tcp",
  "state": "open",
  "scanned_at": "2024-01-15T12:00:00Z"
}
```
