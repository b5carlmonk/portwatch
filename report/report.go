package report

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/portwatch/scanner"
)

// Reporter formats and writes scan results and diffs.
type Reporter struct {
	out io.Writer
}

// New creates a Reporter writing to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{out: w}
}

// PrintResults writes a formatted table of scan results.
func (r *Reporter) PrintResults(results []scanner.Result, scannedAt time.Time) {
	fmt.Fprintf(r.out, "Scan results for %s at %s\n", hostFrom(results), scannedAt.Format(time.RFC3339))
	tw := tabwriter.NewWriter(r.out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "PORT\tPROTO\tSTATE")
	for _, res := range results {
		if res.Open {
			fmt.Fprintf(tw, "%d\t%s\topen\n", res.Port, res.Proto)
		}
	}
	tw.Flush()
}

// PrintDiff writes a human-readable summary of port changes.
func (r *Reporter) PrintDiff(diff scanner.Diff) {
	if len(diff.Opened) == 0 && len(diff.Closed) == 0 {
		fmt.Fprintln(r.out, "No changes detected.")
		return
	}
	for _, res := range diff.Opened {
		fmt.Fprintf(r.out, "[+] Port %d/%s is now OPEN\n", res.Port, res.Proto)
	}
	for _, res := range diff.Closed {
		fmt.Fprintf(r.out, "[-] Port %d/%s is now CLOSED\n", res.Port, res.Proto)
	}
}

func hostFrom(results []scanner.Result) string {
	if len(results) > 0 {
		return results[0].Host
	}
	return "unknown"
}
