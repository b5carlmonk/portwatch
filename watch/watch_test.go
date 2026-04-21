package watch_test

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/alert"
	"github.com/user/portwatch/watch"
)

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestWatcherRunSucceeds(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")

	a := alert.New(os.Stdout)
	cfg := watch.Config{
		Targets:   []string{"127.0.0.1:9"},
		StateFile: statePath,
		TimeoutMs: 200,
	}
	w, err := watch.New(cfg, a)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := w.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if _, err := os.Stat(statePath); err != nil {
		t.Errorf("expected state file to exist: %v", err)
	}
}

func TestWatcherRunPersistsState(t *testing.T) {
	dir := t.TempDir()
	statePath := filepath.Join(dir, "state.json")

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer l.Close()

	a := alert.New(os.Stdout)
	cfg := watch.Config{
		Targets:   []string{l.Addr().String()},
		StateFile: statePath,
		TimeoutMs: 500,
	}
	w, err := watch.New(cfg, a)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Run twice; second run should see no diff (no alerts).
	for i := 0; i < 2; i++ {
		if err := w.Run(context.Background()); err != nil {
			t.Fatalf("Run %d: %v", i, err)
		}
	}
}

func TestWatcherNoTargetsErrors(t *testing.T) {
	a := alert.New(os.Stdout)
	_, err := watch.New(watch.Config{StateFile: "x.json"}, a)
	if err == nil {
		t.Fatal("expected error for empty targets")
	}
}
