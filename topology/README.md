# topology

The `topology` package maintains a live map of hosts to their observed open ports,
giving portwatch a structured view of the current network surface.

## Usage

```go
topo := topology.New()

// Feed scanner results into the topology after every scan cycle.
topo.Update(results)

// Query a specific host.
if node, ok := topo.Get("192.168.1.1"); ok {
    fmt.Printf("%s has %d open ports\n", node.Host, len(node.Ports))
}

// List all known hosts (sorted).
for _, h := range topo.Hosts() {
    fmt.Println(h)
}

// Print a human-readable summary.
fmt.Print(topo.Summary())
```

## Behaviour

- `Update` replaces the port list for each host present in the provided results.
  Hosts not included in the new results are **not** removed — they remain until
  the next update that includes them.
- `Get` is safe for concurrent reads.
- `Hosts` returns host addresses in lexicographic order.
- `Summary` produces a compact, multi-line string suitable for logging or
  terminal output.

## Thread Safety

All exported methods are protected by an internal `sync.RWMutex` and are safe
to call from multiple goroutines.
