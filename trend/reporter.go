package trend

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Reporter prints trend summaries.
type Reporter struct {
	w io.Writer
}

// NewReporter returns a Reporter writing to w. If w is nil, os.Stdout is used.
func NewReporter(w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{w: w}
}

// Print writes a sorted summary of all trend entries.
func (r *Reporter) Print(t *Tracker) {
	entries := t.All()
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})
	fmt.Fprintln(r.w, "--- Port Trend Summary ---")
	for _, e := range entries {
		fmt.Fprintf(r.w, "  %-30s opened=%-4d closed=%-4d last=%s\n",
			e.Key, e.OpenCount, e.CloseCount, e.LastSeen.Format("2006-01-02 15:04:05"))
	}
	if len(entries) == 0 {
		fmt.Fprintln(r.w, "  (no data)")
	}
}
