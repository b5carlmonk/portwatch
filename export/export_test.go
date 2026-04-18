package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/export"
	"github.com/user/portwatch/scanner"
)

var epoch = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "localhost", Port: 80, Protocol: "tcp", State: "open", ScannedAt: epoch},
		{Host: "localhost", Port: 443, Protocol: "tcp", State: "open", ScannedAt: epoch},
	}
}

func TestJSONContainsPorts(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatJSON)
	if err := e.Write(makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if int(records[0]["port"].(float64)) != 80 {
		t.Errorf("expected port 80, got %v", records[0]["port"])
	}
}

func TestJSONContainsHost(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatJSON)
	_ = e.Write(makeResults())
	if !strings.Contains(buf.String(), "localhost") {
		t.Error("expected host 'localhost' in JSON output")
	}
}

func TestCSVHeader(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatCSV)
	if err := e.Write(makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if lines[0] != "host,port,protocol,state,scanned_at" {
		t.Errorf("unexpected header: %s", lines[0])
	}
}

func TestCSVRowCount(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatCSV)
	_ = e.Write(makeResults())
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// header + 2 data rows
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestCSVContainsScannedAt(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.FormatCSV)
	if err := e.Write(makeResults()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Verify the epoch timestamp appears in the CSV output.
	if !strings.Contains(buf.String(), "2024-01-15") {
		t.Error("expected scanned_at date '2024-01-15' in CSV output")
	}
}

func TestUnknownFormatReturnsError(t *testing.T) {
	var buf bytes.Buffer
	e := export.New(&buf, export.Format("xml"))
	if err := e.Write(makeResults()); err == nil {
		t.Error("expected error for unknown format")
	}
}
