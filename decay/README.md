# decay

The `decay` package implements a time-based exponential score decay tracker.

## Overview

Each key (e.g. `"192.168.1.1:443"`) carries a floating-point score that
automatically decreases over time using the formula:

```
score(t) = score₀ × 0.5^(elapsed / halfLife)
```

After one `halfLife` duration the score is halved; after two half-lives it is
quartered, and so on.

## Usage

```go
tr := decay.New(30 * time.Minute)

// Record a risk event for a port.
tr.Add("10.0.0.1:8080", 50)

// Later — score will be lower depending on elapsed time.
current := tr.Get("10.0.0.1:8080")
fmt.Printf("current risk score: %.2f\n", current)

// Manually clear a key.
tr.Reset("10.0.0.1:8080")

// Clear all keys.
tr.Flush()
```

## Integration

Combine with `scorecard` or `classify` to produce time-aware risk scores that
fade when a previously risky port stops appearing in scan results.
