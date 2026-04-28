// Package anomaly detects statistical anomalies in port scan results
// by comparing observed open-port counts against a rolling baseline mean.
package anomaly

import (
	"fmt"
	"math"
	"sync"

	"github.com/user/portwatch/scanner"
)

// Detector tracks rolling statistics and flags scans that deviate
// significantly from the historical mean.
type Detector struct {
	mu        sync.Mutex
	threshold float64 // standard-deviation multiplier (e.g. 2.0)
	samples   []float64
}

// Alert describes a detected anomaly.
type Alert struct {
	Host    string
	Observed int
	Mean    float64
	StdDev  float64
	Message string
}

// New returns a Detector that fires when the observed count is more than
// threshold standard deviations away from the rolling mean.
func New(threshold float64) *Detector {
	if threshold <= 0 {
		threshold = 2.0
	}
	return &Detector{threshold: threshold}
}

// Analyze records the number of open ports in results and returns an Alert
// if the count is anomalous. Returns nil when no anomaly is detected or when
// there is insufficient history (fewer than 3 samples).
func (d *Detector) Analyze(results []scanner.Result) *Alert {
	open := countOpen(results)

	d.mu.Lock()
	defer d.mu.Unlock()

	d.samples = append(d.samples, float64(open))

	if len(d.samples) < 3 {
		return nil
	}

	mean, std := stats(d.samples)
	if std == 0 {
		return nil
	}

	z := math.Abs(float64(open)-mean) / std
	if z < d.threshold {
		return nil
	}

	host := ""
	if len(results) > 0 {
		host = results[0].Host
	}

	return &Alert{
		Host:     host,
		Observed: open,
		Mean:     mean,
		StdDev:   std,
		Message:  fmt.Sprintf("open ports %d deviates %.2f stddev from mean %.2f", open, z, mean),
	}
}

// Reset clears all recorded samples.
func (d *Detector) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.samples = nil
}

func countOpen(results []scanner.Result) int {
	n := 0
	for _, r := range results {
		if r.Open {
			n++
		}
	}
	return n
}

func stats(s []float64) (mean, std float64) {
	for _, v := range s {
		mean += v
	}
	mean /= float64(len(s))
	for _, v := range s {
		d := v - mean
		std += d * d
	}
	std = math.Sqrt(std / float64(len(s)))
	return
}
