package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/ratelimit"
)

func TestAllowFirstEvent(t *testing.T) {
	l := ratelimit.New(time.Second)
	if !l.Allow("tcp:22") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestSuppressesWithinInterval(t *testing.T) {
	l := ratelimit.New(time.Hour)
	l.Allow("tcp:80")
	if l.Allow("tcp:80") {
		t.Fatal("expected second event within interval to be suppressed")
	}
}

func TestAllowsAfterInterval(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("tcp:443")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("tcp:443") {
		t.Fatal("expected event to be allowed after interval expires")
	}
}

func TestResetClearsKey(t *testing.T) {
	l := ratelimit.New(time.Hour)
	l.Allow("tcp:8080")
	l.Reset("tcp:8080")
	if !l.Allow("tcp:8080") {
		t.Fatal("expected event to be allowed after reset")
	}
}

func TestFlushClearsAll(t *testing.T) {
	l := ratelimit.New(time.Hour)
	l.Allow("tcp:22")
	l.Allow("tcp:80")
	l.Flush()
	if !l.Allow("tcp:22") {
		t.Fatal("expected tcp:22 to be allowed after flush")
	}
	if !l.Allow("tcp:80") {
		t.Fatal("expected tcp:80 to be allowed after flush")
	}
}

func TestIndependentKeys(t *testing.T) {
	l := ratelimit.New(time.Hour)
	l.Allow("tcp:22")
	if !l.Allow("tcp:80") {
		t.Fatal("expected different key to be allowed independently")
	}
}
