# scorecard

The `scorecard` package computes a **risk score** for a scanned host based on
the number and type of open ports discovered by portwatch.

## How scoring works

| Condition | Points |
|-----------|--------|
| Each open port | +2 |
| Sensitive port (Telnet, SMB, RDP, …) | +20 – +40 |

### Risk levels

| Score | Level |
|-------|-------|
| 0 – 19 | `low` |
| 20 – 49 | `medium` |
| 50 – 79 | `high` |
| 80+ | `critical` |

## Usage

```go
scorer := scorecard.New()
report := scorer.Score(results)
fmt.Println(report.Level, report.Score)
```

## Reporting

```go
rep := scorecard.NewReporter()
rep.Print(report)
```

Output example:

```
Host     : 192.168.1.10
Score    : 76
Level    : HIGH
Open     : 4 ports
Breakdown:
  - port 23 (tcp) +40
  - port 445 (tcp) +35
```

## History

Each `Scorer` instance keeps an in-memory history of all reports per host:

```go
past := scorer.History("192.168.1.10")
```
