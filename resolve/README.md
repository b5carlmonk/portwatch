# resolve

The `resolve` package enriches scan results with reverse-DNS hostnames.

## Usage

```go
r := resolve.New()
enriched := r.Enrich(results)
```

## Behaviour

- Calls `net.LookupAddr` for each unique IP.
- Results are cached in-memory for the lifetime of the `Resolver`.
- If lookup fails the original IP is kept as the host value.
- `Enrich` returns a new slice and does not mutate the input.

## Testing

Inject a custom lookup function via `NewWithLookup` to avoid real DNS calls in tests.
