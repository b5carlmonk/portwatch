package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the portwatch configuration.
type Config struct {
	Host     string        `json:"host"`
	Ports    []int         `json:"ports"`
	Interval time.Duration `json:"interval"`
	Alerts   AlertConfig   `json:"alerts"`
}

// AlertConfig holds alerting configuration.
type AlertConfig struct {
	Email   string `json:"email,omitempty"`
	Webhook string `json:"webhook,omitempty"`
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		Host:     "localhost",
		Ports:    []int{22, 80, 443, 8080},
		Interval: 60 * time.Second,
	}
}

// Load reads a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := Default()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Save writes the config as JSON to the given path.
func Save(path string, cfg *Config) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(cfg)
}
