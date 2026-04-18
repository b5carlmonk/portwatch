package watchdog_test

import (
	"context"
	"net"
	"testing"
	"os"

	"portwatch/alert"
	"portwatch/config"
	"portwatch/filter"
	"portwatch/history"
	"portwatch/report"
	"portwatch/scanner"
	"portwatch/watchdog"
)

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestRunCycleNoError(t *testing.T) {
	dir := t.TempDir()

	cfg := config.Default()
	cfg.Host = "127.0.0.1"
	cfg.StateFile = dir + "/state.json"
	cfg.HistoryFile = dir + "/history.json"
	cfg.Ports = []int{freePort(t)}

	sc := scanner.New(cfg.Timeout)
	f := filter.New(nil)
	a := alert.New(cfg)
	r := report.New(os.Stdout)
	h := history.New(10)

	wd := watchdog.New(cfg, sc, f, a, r, h)
	if err := wd.Run(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunCyclePersistsState(t *testing.T) {
	dir := t.TempDir()

	cfg := config.Default()
	cfg.Host = "127.0.0.1"
	cfg.StateFile = dir + "/state.json"
	cfg.HistoryFile = dir + "/history.json"
	cfg.Ports = []int{freePort(t)}

	sc := scanner.New(cfg.Timeout)
	f := filter.New(nil)
	a := alert.New(cfg)
	r := report.New(os.Stdout)
	h := history.New(10)

	wd := watchdog.New(cfg, sc, f, a, r, h)
	if err := wd.Run(context.Background()); err != nil {
		t.Fatalf("run error: %v", err)
	}

	if _, err := os.Stat(cfg.StateFile); err != nil {
		t.Fatalf("state file not created: %v", err)
	}
	if _, err := os.Stat(cfg.HistoryFile); err != nil {
		t.Fatalf("history file not created: %v", err)
	}
}
