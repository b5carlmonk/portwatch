# acknowledge

The `acknowledge` package lets operators mark specific port/protocol/host combinations as **acknowledged**, suppressing repeated alerts for known or expected open ports.

## Usage

```go
store := acknowledge.New("ack.json")
_ = store.Load()

k := acknowledge.Key{Host: "localhost", Port: 22, Protocol: "tcp"}
store.Acknowledge(k)

if store.IsAcknowledged(k) {
    // skip alert
}

_ = store.Save()
```

## API

| Function | Description |
|---|---|
| `New(path)` | Create a new store backed by the given JSON file |
| `Acknowledge(Key)` | Mark a key as acknowledged |
| `IsAcknowledged(Key)` | Check whether a key is acknowledged |
| `Revoke(Key)` | Remove an acknowledgement |
| `Save()` | Persist acknowledged keys to disk |
| `Load()` | Load acknowledged keys from disk |

## Persistence

Keys are stored as a JSON array. Missing files are silently ignored on load.
