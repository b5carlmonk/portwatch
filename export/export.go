// Package export writes scan results to various output formats.
package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/scanner"
)

// Format represents an output format.
type Format string

const (
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Exporter writes scan results to an io.Writer.
type Exporter struct {
	w      io.Writer
	format Format
}

// New returns a new Exporter for the given format.
func New(w io.Writer, format Format) *Exporter {
	return &Exporter{w: w, format: format}
}

type jsonRecord struct {
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	State     string    `json:"state"`
	ScannedAt time.Time `json:"scanned_at"`
}

// Write serialises results in the configured format.
func (e *Exporter) Write(results []scanner.Result) error {
	switch e.format {
	case FormatJSON:
		return e.writeJSON(results)
	case FormatCSV:
		return e.writeCSV(results)
	default:
		return fmt.Errorf("export: unknown format %q", e.format)
	}
}

func (e *Exporter) writeJSON(results []scanner.Result) error {
	records := make([]jsonRecord, len(results))
	for i, r := range results {
		records[i] = jsonRecord{
			Host:      r.Host,
			Port:      r.Port,
			Protocol:  r.Protocol,
			State:     r.State,
			ScannedAt: r.ScannedAt,
		}
	}
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(records)
}

func (e *Exporter) writeCSV(results []scanner.Result) error {
	w := csv.NewWriter(e.w)
	if err := w.Write([]string{"host", "port", "protocol", "state", "scanned_at"}); err != nil {
		return err
	}
	for _, r := range results {
		row := []string{
			r.Host,
			fmt.Sprintf("%d", r.Port),
			r.Protocol,
			r.State,
			r.ScannedAt.Format(time.RFC3339),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
