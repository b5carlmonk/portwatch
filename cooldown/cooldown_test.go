package cooldown

import (
	"testing"
	"time"
)

func newFixed(t time.Time) *Cooldown {
	c := New(5 * time.Second)
	c.now = func() time.Time { return t }
	return c
}

func TestAllowFirstEvent(t *testing.T) {
	c := newFixed(time.Now())
	if !c.Allow("port:80") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestSuppressesWithinCooldown(t *testing.T) {
	base := time.Now()
	c := newFixed(base)
	c.Allow("port:80")

	// advance by less than cooldown duration
	c.now = func() time.Time { return base.Add(3 * time.Second) }
	if c.Allow("port:80") {
		t.Fatal("expected call within cooldown to be suppressed")
	}
}

func TestAllowsAfterCooldownExpires(t *testing.T) {
	base := time.Now()
	c := newFixed(base)
	c.Allow("port:80")

	// advance past cooldown duration
	c.now = func() time.Time { return base.Add(6 * time.Second) }
	if !c.Allow("port:80") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestResetClearsKey(t *testing.T) {
	base := time.Now()
	c := newFixed(base)
	c.Allow("port:443")

	c.Reset("port:443")
	if !c.Allow("port:443") {
		t.Fatal("expected allow after reset")
	}
}

func TestFlushClearsAll(t *testing.T) {
	base := time.Now()
	c := newFixed(base)
	c.Allow("port:80")
	c.Allow("port:443")

	c.Flush()
	for _, key := range []string{"port:80", "port:443"} {
		if !c.Allow(key) {
			t.Fatalf("expected allow after flush for key %s", key)
		}
	}
}

func TestRemainingWithinCooldown(t *testing.T) {
	base := time.Now()
	c := newFixed(base)
	c.Allow("port:22")

	c.now = func() time.Time { return base.Add(2 * time.Second) }
	rem := c.Remaining("port:22")
	if rem <= 0 {
		t.Fatalf("expected positive remaining, got %v", rem)
	}
}

func TestRemainingUnknownKey(t *testing.T) {
	c := newFixed(time.Now())
	if r := c.Remaining("unknown"); r != 0 {
		t.Fatalf("expected 0 for unknown key, got %v", r)
	}
}
