# priority

The `priority` package assigns severity levels to open port scan results
based on user-defined rules.

## Levels

| Level    | Value |
|----------|-------|
| Low      | 0     |
| Medium   | 1     |
| High     | 2     |
| Critical | 3     |

## Usage

```go
scorer := priority.New([]priority.Rule{
    {Port: 22,  Proto: "tcp", Level: priority.Critical},
    {Port: 443, Proto: "tcp", Level: priority.High},
    {Port: 53,  Proto: "",    Level: priority.Medium}, // any protocol
}, priority.Low) // fallback for unmatched ports

lvl := scorer.Score(result)          // single result
map  := scorer.ScoreAll(results)     // all results -> map["port/proto"]Level
```

## Reporter

```go
r := priority.NewReporter()
r.Print(scorer.ScoreAll(results))
```

Output:

```
PORT/PROTO       PRIORITY
--------------------------
22/tcp           CRITICAL
443/tcp          HIGH
8080/tcp         LOW
```

## Rules

Each `Rule` specifies:
- `Port` – the port number to match.
- `Proto` – `"tcp"`, `"udp"`, or `""` to match any protocol.
- `Level` – the severity assigned when the rule matches.

Rules are evaluated in order; the first match wins.
