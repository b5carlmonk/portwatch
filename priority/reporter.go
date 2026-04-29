package priority

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Reporter prints a priority score table to a writer.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter that writes to stdout.
func NewReporter() *Reporter { return &Reporter{w: os.Stdout} }

// NewReporterWithWriter returns a Reporter that writes to w.
func NewReporterWithWriter(w io.Writer) *Reporter { return &Reporter{w: w} }

// Print writes a formatted table of port priorities to the writer.
func (r *Reporter) Print(scores map[string]Level) {
	if len(scores) == 0 {
		return
	}

	keys := make([]string, 0, len(scores))
	for k := range scores {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintln(r.w, "PORT/PROTO       PRIORITY")
	fmt.Fprintln(r.w, strings.Repeat("-", 26))
	for _, k := range keys {
		fmt.Fprintf(r.w, "%-16s %s\n", k, scores[k])
	}
}
