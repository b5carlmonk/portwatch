package anomaly

import (
	"fmt"
	"io"
	"os"
)

// Reporter writes anomaly alerts to an io.Writer.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter that writes to stdout.
func NewReporter() *Reporter {
	return &Reporter{w: os.Stdout}
}

// NewReporterWithWriter returns a Reporter that writes to w.
func NewReporterWithWriter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Print formats and writes the alert to the underlying writer.
// If alert is nil, nothing is written.
func (r *Reporter) Print(alert *Alert) {
	if alert == nil {
		return
	}
	fmt.Fprintf(r.w, "[ANOMALY] host=%s observed=%d mean=%.2f stddev=%.2f\n",
		alert.Host, alert.Observed, alert.Mean, alert.StdDev)
	fmt.Fprintf(r.w, "          %s\n", alert.Message)
}
