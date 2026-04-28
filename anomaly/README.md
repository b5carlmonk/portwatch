# anomaly

The `anomaly` package detects statistical anomalies in port scan results by
comparing the number of open ports observed in each scan against a rolling
historical mean and standard deviation.

## How it works

1. Each call to `Analyze` records the open-port count for the current scan.
2. Once at least **3 samples** have been collected, the detector computes the
   population mean and standard deviation.
3. If the current observation deviates from the mean by more than `threshold`
   standard deviations (default **2.0**), an `Alert` is returned.

## Usage

```go
det := anomaly.New(2.0)

results := scanner.Scan(target)
if alert := det.Analyze(results); alert != nil {
    rep := anomaly.NewReporter()
    rep.Print(alert)
}
```

## Alert fields

| Field      | Description                                      |
|------------|--------------------------------------------------|
| `Host`     | Host extracted from the first scan result        |
| `Observed` | Number of open ports in the current scan         |
| `Mean`     | Rolling mean of historical open-port counts      |
| `StdDev`   | Rolling standard deviation                       |
| `Message`  | Human-readable summary of the anomaly            |

## Resetting

Call `Reset()` to discard all historical samples, for example when switching
between scan targets.
