package throttle

import (
	"testing"
	"time"
)

func TestAllowFirstCall(t *testing.T) {
	th := New(time.Second)
	if !th.Allow("host1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestSuppressesWithinInterval(t *testing.T) {
	now := time.Now()
	th := New(time.Minute)
	th.now = func() time.Time { return now }
	th.Allow("host1")
	if th.Allow("host1") {
		t.Fatal("expected second call within interval to be suppressed")
	}
}

func TestAllowsAfterInterval(t *testing.T) {
	now := time.Now()
	th := New(time.Minute)
	th.now = func() time.Time { return now }
	th.Allow("host1")
	th.now = func() time.Time { return now.Add(2 * time.Minute) }
	if !th.Allow("host1") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestResetClearsKey(t *testing.T) {
	now := time.Now()
	th := New(time.Minute)
	th.now = func() time.Time { return now }
	th.Allow("host1")
	th.Reset("host1")
	if !th.Allow("host1") {
		t.Fatal("expected allow after reset")
	}
}

func TestFlushClearsAll(t *testing.T) {
	now := time.Now()
	th := New(time.Minute)
	th.now = func() time.Time { return now }
	th.Allow("host1")
	th.Allow("host2")
	th.Flush()
	if !th.Allow("host1") || !th.Allow("host2") {
		t.Fatal("expected all keys to be cleared after flush")
	}
}

func TestRemainingIsZeroWhenNotSeen(t *testing.T) {
	th := New(time.Minute)
	if r := th.Remaining("host1"); r != 0 {
		t.Fatalf("expected 0 remaining for unseen key, got %v", r)
	}
}

func TestRemainingReturnsPositiveDuration(t *testing.T) {
	now := time.Now()
	th := New(time.Minute)
	th.now = func() time.Time { return now }
	th.Allow("host1")
	th.now = func() time.Time { return now.Add(10 * time.Second) }
	r := th.Remaining("host1")
	if r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}
