package decay

import (
	"testing"
	"time"
)

func TestAddStoresScore(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	tr := newWithClock(time.Minute, clock)

	tr.Add("host:80", 10)
	got := tr.Get("host:80")
	if got != 10 {
		t.Fatalf("expected 10, got %f", got)
	}
}

func TestGetMissingKeyReturnsZero(t *testing.T) {
	tr := New(time.Minute)
	if got := tr.Get("missing"); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestDecayReducesScoreAfterHalfLife(t *testing.T) {
	start := time.Now()
	clock := func() time.Time { return start }
	tr := newWithClock(time.Hour, clock)

	tr.Add("k", 100)

	// Advance clock by exactly one half-life.
	clock = func() time.Time { return start.Add(time.Hour) }
	tr.clock = clock

	got := tr.Get("k")
	const want = 50.0
	const epsilon = 0.001
	if got < want-epsilon || got > want+epsilon {
		t.Fatalf("expected ~%f after one half-life, got %f", want, got)
	}
}

func TestDecayReducesScoreAfterTwoHalfLives(t *testing.T) {
	start := time.Now()
	clock := func() time.Time { return start }
	tr := newWithClock(time.Hour, clock)

	tr.Add("k", 100)

	clock = func() time.Time { return start.Add(2 * time.Hour) }
	tr.clock = clock

	got := tr.Get("k")
	const want = 25.0
	const epsilon = 0.001
	if got < want-epsilon || got > want+epsilon {
		t.Fatalf("expected ~%f after two half-lives, got %f", want, got)
	}
}

func TestAddAccumulatesAfterDecay(t *testing.T) {
	start := time.Now()
	clock := func() time.Time { return start }
	tr := newWithClock(time.Hour, clock)

	tr.Add("k", 80)

	// Advance one half-life, score should be ~40, then add 20 → ~60.
	clock = func() time.Time { return start.Add(time.Hour) }
	tr.clock = clock
	tr.Add("k", 20)

	got := tr.Get("k")
	const want = 60.0
	const epsilon = 0.01
	if got < want-epsilon || got > want+epsilon {
		t.Fatalf("expected ~%f, got %f", want, got)
	}
}

func TestResetClearsKey(t *testing.T) {
	tr := New(time.Minute)
	tr.Add("k", 50)
	tr.Reset("k")
	if got := tr.Get("k"); got != 0 {
		t.Fatalf("expected 0 after reset, got %f", got)
	}
}

func TestFlushClearsAll(t *testing.T) {
	tr := New(time.Minute)
	tr.Add("a", 10)
	tr.Add("b", 20)
	tr.Flush()
	if got := tr.Get("a"); got != 0 {
		t.Fatalf("expected 0 after flush, got %f", got)
	}
	if got := tr.Get("b"); got != 0 {
		t.Fatalf("expected 0 after flush, got %f", got)
	}
}
