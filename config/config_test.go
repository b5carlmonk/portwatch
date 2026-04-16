package config

import (
	"os"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", cfg.Host)
	}
	if len(cfg.Ports) == 0 {
		t.Error("expected default ports to be non-empty")
	}
	if cfg.Interval != 60*time.Second {
		t.Errorf("expected 60s interval, got %v", cfg.Interval)
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-config-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	cfg := &Config{
		Host:     "192.168.1.1",
		Ports:    []int{22, 443, 8443},
		Interval: 30 * time.Second,
		Alerts:   AlertConfig{Webhook: "http://example.com/hook"},
	}

	if err := Save(tmp.Name(), cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Host != cfg.Host {
		t.Errorf("host mismatch: got %s, want %s", loaded.Host, cfg.Host)
	}
	if len(loaded.Ports) != len(cfg.Ports) {
		t.Errorf("ports length mismatch: got %d, want %d", len(loaded.Ports), len(cfg.Ports))
	}
	if loaded.Interval != cfg.Interval {
		t.Errorf("interval mismatch: got %v, want %v", loaded.Interval, cfg.Interval)
	}
	if loaded.Alerts.Webhook != cfg.Alerts.Webhook {
		t.Errorf("webhook mismatch: got %s, want %s", loaded.Alerts.Webhook, cfg.Alerts.Webhook)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.json")
	if err == nil {
		t.Error("expected error loading missing file, got nil")
	}
}
