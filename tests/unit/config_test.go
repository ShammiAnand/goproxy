package unit

import (
	"github.com/shammianand/goproxy/internal/config"
	"log/slog"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	content := []byte(`
listen_addr: ":8080"
target_addr: "http://example.com"
log_level: "debug"
`)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Test loading the config
	cfg, err := config.Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.ListenAddr != ":8080" {
		t.Errorf("Expected ListenAddr to be ':8080', got '%s'", cfg.ListenAddr)
	}
	if cfg.TargetAddr != "http://example.com" {
		t.Errorf("Expected TargetAddr to be 'http://example.com', got '%s'", cfg.TargetAddr)
	}
	if cfg.GetLogLevel() != slog.LevelDebug {
		t.Errorf("Expected LogLevel to be debug, got '%s'", cfg.LogLevel)
	}
}
