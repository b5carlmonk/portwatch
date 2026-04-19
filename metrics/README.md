# metrics

The `metrics` package tracks runtime counters for portwatch scan cycles.

## Overview

A `Tracker` accumulates statistics across scan cycles and can print a summary line at any time.

## Counters

| Field | Description |
|-------|-------------|
| `Scans` | Total number of scan cycles completed |
| `OpenPorts` | Open port count from the most recent scan |
| `Changes` | Cumulative number of port state changes detected |
| `Errors` | Total scan cycles that returned an error |
| `LastScan` | Timestamp of the most recent scan |

## Usage

```go
tr := metrics.New()

// after each scan cycle:
tr.RecordScan(len(results), len(diff.Opened)+len(diff.Closed), err)

// print summary:
tr.Print()
// scans=12 open=5 changes=3 errors=0 last=2024-01-15T10:30:00Z

// or read counters programmatically:
s := tr.Snapshot()
fmt.Println(s.Scans, s.Errors)
```
