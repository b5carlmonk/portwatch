package scorecard

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// Reporter prints scorecard reports to a writer.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter writing to stdout.
func NewReporter() *Reporter { return &Reporter{w: os.Stdout} }

// NewReporterWithWriter returns a Reporter writing to w.
func NewReporterWithWriter(w io.Writer) *Reporter { return &Reporter{w: w} }

// Print writes a formatted report to the underlying writer.
func (r *Reporter) Print(rep Report) {
	fmt.Fprintf(r.w, "Host     : %s\n", rep.Host)
	fmt.Fprintf(r.w, "Score    : %d\n", rep.Score)
	fmt.Fprintf(r.w, "Level    : %s\n", strings.ToUpper(rep.Level))
	fmt.Fprintf(r.w, "Open     : %d ports\n", rep.OpenPorts)
	if len(rep.Breakdown) > 0 {
		fmt.Fprintln(r.w, "Breakdown:")
		for _, b := range rep.Breakdown {
			fmt.Fprintf(r.w, "  - %s\n", b)
		}
	}
}
