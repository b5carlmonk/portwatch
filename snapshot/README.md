# snapshot

The `snapshot` package provides named, persistent captures of port scan results.

## Usage

```go
store := snapshot.New("snapshots.json")
_ = store.Load()

// Save current scan as a named snapshot
_ = store.Add("before-deploy", results)

// Retrieve later
entry, ok := store.Get("before-deploy")

// List all labels
fmt.Println(store.Labels())

// Remove when no longer needed
_ = store.Delete("before-deploy")
```

## Entry fields

| Field | Description |
|-------|-------------|
| `Label` | Human-readable name |
| `CreatedAt` | UTC timestamp of capture |
| `Results` | Slice of `scanner.Result` |

Snapshots are persisted as JSON and survive process restarts.
