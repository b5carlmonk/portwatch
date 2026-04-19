# healthcheck

The `healthcheck` package provides active TCP/UDP probing for known ports to verify they are reachable and measure latency.

## Usage

```go
c := healthcheck.New(2 * time.Second)

// Probe a single port
result := c.Probe("192.168.1.1", 80, "tcp")
fmt.Println(result.Alive, result.Latency)

// Probe multiple targets
targets := []healthcheck.Target{
    {Port: 80, Proto: "tcp"},
    {Port: 443, Proto: "tcp"},
}
results := c.ProbeAll("192.168.1.1", targets)
for _, r := range results {
    fmt.Printf("port %d alive=%v latency=%v\n", r.Port, r.Alive, r.Latency)
}
```

## Result Fields

| Field   | Description                        |
|---------|------------------------------------|
| Host    | Target host                        |
| Port    | Target port                        |
| Proto   | Protocol used (tcp/udp)            |
| Alive   | Whether the port responded         |
| Latency | Round-trip dial time               |
| Err     | Error if probe failed, else nil    |
