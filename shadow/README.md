# shadow

The `shadow` package detects **shadow services** — ports that are open on a
host but were not present in a known-good baseline snapshot.

## Concepts

| Term | Meaning |
|------|---------|
| Baseline | A set of `scanner.Result` values captured at a trusted point in time |
| Shadow port | An open port found during a live scan that is absent from the baseline |
| Entry | A single shadow finding, including host, port, protocol, and reason |

## Usage

```go
// Capture a baseline (e.g. from state or sampler)
baseline := []scanner.Result{ /* trusted results */ }

d := shadow.New(baseline)

// Later, compare against a fresh scan
current := scanner.Scan(targets, opts)
findings := d.Detect(current)

for _, f := range findings {
    fmt.Printf("SHADOW %s:%d/%s — %s\n", f.Host, f.Port, f.Protocol, f.Reason)
}
```

## Runtime updates

Use `Allow` to whitelist a port without rebuilding the detector:

```go
d.Allow(scanner.Result{Host: "localhost", Port: 8080, Protocol: "tcp", Open: true})
```

## Notes

- Only **open** ports are considered; closed ports in both baseline and current
  scan are ignored.
- The detector is safe for concurrent use.
