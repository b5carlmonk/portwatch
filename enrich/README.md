# enrich

The `enrich` package attaches additional metadata to scan results before they are processed by downstream components.

## Overview

An `Enricher` accepts one or more `Provider` functions. Each provider receives a `scanner.Result` and returns a key/value pair that is stored in the result's `Meta.Extra` map.

Two built-in providers are included:

- `PortLabel` — adds a `port_label` entry (e.g. `port-80`)
- `ProtoLabel` — adds a `proto_label` entry (e.g. `tcp`)

## Usage

```go
e := enrich.New(enrich.PortLabel, enrich.ProtoLabel)
enriched := e.Enrich(results)
for _, r := range enriched {
    fmt.Println(r.Port, r.Meta.Extra["port_label"], r.Meta.ScannedAt)
}
```

## Custom Providers

```go
secureTag := func(r scanner.Result) (string, string) {
    if r.Port == 443 {
        return "secure", "true"
    }
    return "", ""
}
e := enrich.New(secureTag)
```

Providers that return an empty key are silently ignored.
