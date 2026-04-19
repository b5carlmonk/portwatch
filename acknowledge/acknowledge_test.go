package acknowledge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/acknowledge"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "ack.json")
}

func TestAcknowledgeAndCheck(t *testing.T) {
	s := acknowledge.New(tmpPath(t))
	k := acknowledge.Key{Host: "localhost", Port: 22, Protocol: "tcp"}
	if s.IsAcknowledged(k) {
		t.Fatal("expected not acknowledged")
	}
	s.Acknowledge(k)
	if !s.IsAcknowledged(k) {
		t.Fatal("expected acknowledged")
	}
}

func TestRevoke(t *testing.T) {
	s := acknowledge.New(tmpPath(t))
	k := acknowledge.Key{Host: "localhost", Port: 80, Protocol: "tcp"}
	s.Acknowledge(k)
	s.Revoke(k)
	if s.IsAcknowledged(k) {
		t.Fatal("expected revoked")
	}
}

func TestSaveAndLoad(t *testing.T) {
	p := tmpPath(t)
	s := acknowledge.New(p)
	k := acknowledge.Key{Host: "host1", Port: 443, Protocol: "tcp"}
	s.Acknowledge(k)
	if err := s.Save(); err != nil {
		t.Fatalf("save: %v", err)
	}
	s2 := acknowledge.New(p)
	if err := s2.Load(); err != nil {
		t.Fatalf("load: %v", err)
	}
	if !s2.IsAcknowledged(k) {
		t.Fatal("expected key to survive save/load")
	}
}

func TestLoadMissingFile(t *testing.T) {
	s := acknowledge.New("/tmp/portwatch_ack_missing_xyz.json")
	if err := s.Load(); err != nil {
		t.Fatalf("expected no error on missing file, got %v", err)
	}
}

func TestLoadDoesNotDuplicateKeys(t *testing.T) {
	p := tmpPath(t)
	s := acknowledge.New(p)
	k := acknowledge.Key{Host: "h", Port: 9090, Protocol: "udp"}
	s.Acknowledge(k)
	_ = s.Save()
	_ = s.Load() // load into already-populated store
	// Revoke once — should be gone
	s.Revoke(k)
	if s.IsAcknowledged(k) {
		t.Fatal("duplicate key caused ghost entry")
	}
}

func TestSaveCreatesFile(t *testing.T) {
	p := tmpPath(t)
	s := acknowledge.New(p)
	s.Acknowledge(acknowledge.Key{Host: "x", Port: 1, Protocol: "tcp"})
	if err := s.Save(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}
