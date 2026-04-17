# report

The `report` package formats and displays port scan results and change diffs to the terminal.

## Usage

```go
r := report.New(nil) // writes to stdout

// Print all open ports from a scan
r.PrintResults(results, time.Now())

// Print what changed since the last scan
r.PrintDiff(diff)
```

## Output examples

### PrintResults

```
Scan results for localhost at 2024-01-15T10:00:00Z
PORT  PROTO  STATE
22    tcp    open
80    tcp    open
443   tcp    open
```

### PrintDiff

```
[+] Port 8080/tcp is now OPEN
[-] Port 3306/tcp is now CLOSED
```
