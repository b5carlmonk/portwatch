# classify

The `classify` package maps open ports to human-readable service names.

## Overview

A `Classifier` holds a lookup table of `port/proto → service name`. It ships
with a built-in table of well-known services (SSH, HTTP, MySQL, Redis, etc.) and
accepts an optional custom map that takes precedence over the defaults.

## Usage

```go
c := classify.New(map[string]string{
    "9200/tcp": "Elasticsearch",
})

// Label a single port
name := c.Label(80, "tcp") // "HTTP"

// Enrich a slice of scanner results in place
results = c.Enrich(results)
for _, r := range results {
    fmt.Printf("%d/%s → %s\n", r.Port, r.Proto, r.Service)
}
```

## Built-in mappings

| Port | Proto | Service |
|------|-------|---------|
| 22 | tcp | SSH |
| 80 | tcp | HTTP |
| 443 | tcp | HTTPS |
| 53 | tcp/udp | DNS |
| 3306 | tcp | MySQL |
| 5432 | tcp | PostgreSQL |
| 6379 | tcp | Redis |
| 27017 | tcp | MongoDB |

Unrecognised ports are labelled `"unknown"`.
