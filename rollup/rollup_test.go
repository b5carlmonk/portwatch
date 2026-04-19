package rollup

import (
	"testing"
	"time"

	"github.com/user/portwatch/scanner"
)

func makeDiff(opened, closed []string) scanner.Diff {
	d := scanner.Diff{}
	for _, p := range opened {
		d.Opened = append(d.Opened, scanner.Result{Port: p})
	}
	for _, p := range closed {
		d.Closed = append(d.Closed, scanner.Result{Port: p})
	}
	return d
}

func TestAddIncreasesEntries(t *testing.T) {
	r := New(time.Minute)
	r.Add(makeDiff([]string{"80"}, nil))
	r.Add(makeDiff([]string{"443"}, nil))
	if got := len(r.Entries()); got != 2 {
		t.Fatalf("expected 2 entries, got %d", got)
	}
}

func TestSummaryCountsCorrectly(t *testing.T) {
	r := New(time.Minute)
	r.Add(makeDiff([]string{"80", "443"}, []string{"22"}))
	r.Add(makeDiff([]string{"8080"}, []string{"3306", "5432"}))
	opened, closed := r.Summary()
	if opened != 3 {
		t.Errorf("expected 3 opened, got %d", opened)
	}
	if closed != 3 {
		t.Errorf("expected 3 closed, got %d", closed)
	}
}

func TestWindowEvictsOldEntries(t *testing.T) {
	r := New(50 * time.Millisecond)
	r.Add(makeDiff([]string{"80"}, nil))
	time.Sleep(80 * time.Millisecond)
	r.Add(makeDiff([]string{"443"}, nil))
	if got := len(r.Entries()); got != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", got)
	}
}

func TestFlushClearsEntries(t *testing.T) {
	r := New(time.Minute)
	r.Add(makeDiff([]string{"80"}, nil))
	r.Flush()
	if got := len(r.Entries()); got != 0 {
		t.Fatalf("expected 0 entries after flush, got %d", got)
	}
}

func TestEmptyRollupSummary(t *testing.T) {
	r := New(time.Minute)
	opened, closed := r.Summary()
	if opened != 0 || closed != 0 {
		t.Errorf("expected 0,0 got %d,%d", opened, closed)
	}
}
