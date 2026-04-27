package window

import (
	"testing"
	"time"
)

func TestAddIncreasesTotal(t *testing.T) {
	w := New(time.Minute)
	w.Add(3)
	w.Add(7)
	if got := w.Total(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestTotalExcludesExpiredEntries(t *testing.T) {
	now := time.Now()
	w := newWithClock(time.Second, func() time.Time { return now })
	w.Add(5)
	// advance clock past the window
	now = now.Add(2 * time.Second)
	w.Add(2)
	if got := w.Total(); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestLenCountsActiveEntries(t *testing.T) {
	now := time.Now()
	w := newWithClock(time.Second, func() time.Time { return now })
	w.Add(1)
	w.Add(1)
	if w.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", w.Len())
	}
	now = now.Add(2 * time.Second)
	w.Add(1)
	if w.Len() != 1 {
		t.Fatalf("expected 1 entry after eviction, got %d", w.Len())
	}
}

func TestResetClearsEntries(t *testing.T) {
	w := New(time.Minute)
	w.Add(10)
	w.Reset()
	if got := w.Total(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if w.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", w.Len())
	}
}

func TestEmptyWindowReturnsZero(t *testing.T) {
	w := New(time.Minute)
	if got := w.Total(); got != 0 {
		t.Fatalf("expected 0 for empty window, got %d", got)
	}
}

func TestEntriesWithinWindowRetained(t *testing.T) {
	now := time.Now()
	w := newWithClock(10*time.Second, func() time.Time { return now })
	w.Add(4)
	now = now.Add(5 * time.Second)
	w.Add(6)
	// both entries still within 10s window
	if got := w.Total(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}
