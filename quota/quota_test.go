package quota

import (
	"testing"
	"time"
)

func TestAllowFirstScan(t *testing.T) {
	q := New(3, time.Minute)
	ok, err := q.Allow("host-a")
	if !ok || err != nil {
		t.Fatalf("expected first scan to be allowed, got ok=%v err=%v", ok, err)
	}
}

func TestExceedsQuota(t *testing.T) {
	q := New(2, time.Minute)
	q.Allow("host-a") //nolint
	q.Allow("host-a") //nolint
	ok, err := q.Allow("host-a")
	if ok {
		t.Fatal("expected quota to be exceeded")
	}
	if err == nil {
		t.Fatal("expected error when quota exceeded")
	}
}

func TestAllowsAfterWindowExpires(t *testing.T) {
	q := New(1, 50*time.Millisecond)
	q.Allow("host-a") //nolint

	time.Sleep(60 * time.Millisecond)

	ok, err := q.Allow("host-a")
	if !ok || err != nil {
		t.Fatalf("expected scan to be allowed after window expired, got ok=%v err=%v", ok, err)
	}
}

func TestResetClearsKey(t *testing.T) {
	q := New(1, time.Minute)
	q.Allow("host-a") //nolint
	q.Reset("host-a")

	ok, err := q.Allow("host-a")
	if !ok || err != nil {
		t.Fatalf("expected scan to be allowed after reset, got ok=%v err=%v", ok, err)
	}
}

func TestFlushClearsAll(t *testing.T) {
	q := New(1, time.Minute)
	q.Allow("host-a") //nolint
	q.Allow("host-b") //nolint
	q.Flush()

	for _, key := range []string{"host-a", "host-b"} {
		ok, err := q.Allow(key)
		if !ok || err != nil {
			t.Fatalf("expected %q to be allowed after flush", key)
		}
	}
}

func TestRemainingDecrements(t *testing.T) {
	q := New(3, time.Minute)
	if r := q.Remaining("host-a"); r != 3 {
		t.Fatalf("expected 3 remaining, got %d", r)
	}
	q.Allow("host-a") //nolint
	if r := q.Remaining("host-a"); r != 2 {
		t.Fatalf("expected 2 remaining, got %d", r)
	}
}

func TestRemainingZeroWhenExceeded(t *testing.T) {
	q := New(1, time.Minute)
	q.Allow("host-a") //nolint
	if r := q.Remaining("host-a"); r != 0 {
		t.Fatalf("expected 0 remaining, got %d", r)
	}
}

func TestIndependentKeys(t *testing.T) {
	q := New(1, time.Minute)
	q.Allow("host-a") //nolint

	ok, err := q.Allow("host-b")
	if !ok || err != nil {
		t.Fatal("expected host-b to be independent of host-a quota")
	}
}
