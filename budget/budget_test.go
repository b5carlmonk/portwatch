package budget

import (
	"testing"
	"time"
)

func TestAllowWithinBudget(t *testing.T) {
	b := New(10, time.Minute)
	if err := b.Allow(5); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllowExceedsBudget(t *testing.T) {
	b := New(10, time.Minute)
	if err := b.Allow(10); err != nil {
		t.Fatalf("expected first allow to succeed: %v", err)
	}
	if err := b.Allow(1); err != ErrBudgetExceeded {
		t.Fatalf("expected ErrBudgetExceeded, got %v", err)
	}
}

func TestAllowExactLimit(t *testing.T) {
	b := New(5, time.Minute)
	if err := b.Allow(5); err != nil {
		t.Fatalf("expected allow at exact limit: %v", err)
	}
}

func TestRemainingDecreasesAfterAllow(t *testing.T) {
	b := New(10, time.Minute)
	_ = b.Allow(3)
	if got := b.Remaining(); got != 7 {
		t.Fatalf("expected 7 remaining, got %d", got)
	}
}

func TestRemainingIsFullWhenEmpty(t *testing.T) {
	b := New(10, time.Minute)
	if got := b.Remaining(); got != 10 {
		t.Fatalf("expected 10 remaining, got %d", got)
	}
}

func TestResetClearsUsage(t *testing.T) {
	b := New(5, time.Minute)
	_ = b.Allow(5)
	b.Reset()
	if got := b.Remaining(); got != 5 {
		t.Fatalf("expected 5 after reset, got %d", got)
	}
}

func TestWindowEvictsOldEntries(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }

	b := newWithClock(10, time.Minute, clock)
	_ = b.Allow(8)

	// Advance time beyond the window so the old entry is evicted.
	now = now.Add(2 * time.Minute)

	if err := b.Allow(8); err != nil {
		t.Fatalf("expected allow after window expiry: %v", err)
	}
}

func TestAllowAccumulatesAcrossMultipleCalls(t *testing.T) {
	b := New(10, time.Minute)
	for i := 0; i < 10; i++ {
		if err := b.Allow(1); err != nil {
			t.Fatalf("call %d failed: %v", i, err)
		}
	}
	if err := b.Allow(1); err != ErrBudgetExceeded {
		t.Fatalf("expected ErrBudgetExceeded after 10 calls, got %v", err)
	}
}
