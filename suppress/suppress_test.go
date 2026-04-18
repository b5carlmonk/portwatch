package suppress_test

import (
	"testing"

	"github.com/user/portwatch/suppress"
)

func TestBelowThresholdNotConfirmed(t *testing.T) {
	s := suppress.New(3)
	if s.Record("tcp:80:opened") {
		t.Fatal("expected false on first record")
	}
	if s.Record("tcp:80:opened") {
		t.Fatal("expected false on second record")
	}
}

func TestAtThresholdConfirmed(t *testing.T) {
	s := suppress.New(3)
	s.Record("tcp:80:opened")
	s.Record("tcp:80:opened")
	if !s.Record("tcp:80:opened") {
		t.Fatal("expected true at threshold")
	}
}

func TestAboveThresholdNotFiredAgain(t *testing.T) {
	s := suppress.New(2)
	s.Record("key")
	s.Record("key") // confirmed
	if s.Record("key") {
		t.Fatal("should not fire again after threshold")
	}
}

func TestResetClearsCount(t *testing.T) {
	s := suppress.New(2)
	s.Record("key")
	s.Reset("key")
	if s.Count("key") != 0 {
		t.Fatalf("expected 0 after reset, got %d", s.Count("key"))
	}
}

func TestFlushClearsAll(t *testing.T) {
	s := suppress.New(2)
	s.Record("a")
	s.Record("b")
	s.Flush()
	if s.Count("a") != 0 || s.Count("b") != 0 {
		t.Fatal("expected all counts cleared after flush")
	}
}

func TestThresholdOneAlwaysConfirms(t *testing.T) {
	s := suppress.New(1)
	if !s.Record("instant") {
		t.Fatal("threshold=1 should confirm on first record")
	}
}

func TestIndependentKeys(t *testing.T) {
	s := suppress.New(2)
	s.Record("a")
	s.Record("b")
	s.Record("b")
	if s.Count("a") != 1 {
		t.Fatalf("key a should have count 1, got %d", s.Count("a"))
	}
	if s.Count("b") != 2 {
		t.Fatalf("key b should have count 2, got %d", s.Count("b"))
	}
}
