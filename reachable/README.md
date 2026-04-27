# reachable

The `reachable` package provides lightweight host reachability probing before
a full port scan is attempted. It dials a configurable probe port with a short
timeout to determine whether a target host is accessible on the network.

## Usage

```go
checker := reachable.New(2*time.Second, "80")

// Probe a single host
result := checker.Probe(ctx, "192.168.1.1")
if !result.Reachable {
    log.Printf("host unreachable: %v", result.Err)
}

// Probe multiple hosts concurrently
results := checker.ProbeAll(ctx, []string{"192.168.1.1", "192.168.1.2"})
for _, r := range results {
    fmt.Printf("%s reachable=%v latency=%v\n", r.Host, r.Reachable, r.Latency)
}
```

## Caching

Results are cached for 30 seconds by default to avoid redundant network calls
within a single scan cycle. Call `Flush()` to clear the cache between cycles
or when you need a fresh probe.

## Fields

| Field | Description |
|-------|-------------|
| `Host` | The probed hostname or IP |
| `Reachable` | Whether the TCP dial succeeded |
| `Latency` | Round-trip time for the dial |
| `Err` | Non-nil if the probe failed |
