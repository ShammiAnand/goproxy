package unit

import (
	"bytes"
	"github.com/shammianand/goproxy/internal/config"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	content := []byte(`
    server:
      listen_addr: ":8080"
      read_timeout: 5
      write_timeout: 10
      idle_timeout: 120
    proxy:
      target_addr: "http://localhost:8000"
      max_idle_conns: 100
      dial_timeout: 10
    load_balancing:
      enabled: false
      algorithm: "round_robin"
      backends: []
    tls:
      enabled: false
      cert_file: ""
      key_file: ""
    logging:
      level: "debug"
      format: "json"
    metrics:
      enabled: false
      address: ":9090"
    rate_limiting:
      enabled: false
      requests_per_second: 100
      burst: 50
    caching:
      enabled: false
      default_ttl: 300
      max_size_mb: 100
  `)
	tmpfile, err := os.CreateTemp("", "config*.yaml")
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

	// Check server settings
	if cfg.Server.ListenAddr != ":8080" {
		t.Errorf("Expected ListenAddr to be ':8080', got '%s'", cfg.Server.ListenAddr)
	}
	if cfg.GetServerReadTimeout() != 5*time.Second {
		t.Errorf("Expected ReadTimeout to be 5s, got '%s'", cfg.GetServerReadTimeout())
	}
	if cfg.GetServerWriteTimeout() != 10*time.Second {
		t.Errorf("Expected WriteTimeout to be 10s, got '%s'", cfg.GetServerWriteTimeout())
	}
	if cfg.GetServerIdleTimeout() != 120*time.Second {
		t.Errorf("Expected IdleTimeout to be 120s, got '%s'", cfg.GetServerIdleTimeout())
	}

	// Check proxy settings
	if cfg.Proxy.TargetAddr != "http://localhost:8000" {
		t.Errorf("Expected TargetAddr to be 'http://localhost:8000', got '%s'", cfg.Proxy.TargetAddr)
	}
	if cfg.Proxy.MaxIdleConns != 100 {
		t.Errorf("Expected MaxIdleConns to be 100, got '%d'", cfg.Proxy.MaxIdleConns)
	}
	if cfg.GetProxyDialTimeout() != 10*time.Second {
		t.Errorf("Expected DialTimeout to be 10s, got '%s'", cfg.GetProxyDialTimeout())
	}

	// Check load balancing settings
	if cfg.LoadBalancing.Enabled != false {
		t.Errorf("Expected LoadBalancing.Enabled to be false, got '%v'", cfg.LoadBalancing.Enabled)
	}
	if cfg.LoadBalancing.Algorithm != "round_robin" {
		t.Errorf("Expected LoadBalancing.Algorithm to be 'round_robin', got '%s'", cfg.LoadBalancing.Algorithm)
	}

	// Check TLS settings
	if cfg.TLS.Enabled != false {
		t.Errorf("Expected TLS.Enabled to be false, got '%v'", cfg.TLS.Enabled)
	}

	// Check logging settings
	if cfg.GetLogLevel() != slog.LevelDebug {
		t.Errorf("Expected LogLevel to be debug, got '%s'", cfg.Logging.Level)
	}

	// Check metrics settings
	if cfg.Metrics.Enabled != false {
		t.Errorf("Expected Metrics.Enabled to be false, got '%v'", cfg.Metrics.Enabled)
	}
	if cfg.Metrics.Address != ":9090" {
		t.Errorf("Expected Metrics.Address to be ':9090', got '%s'", cfg.Metrics.Address)
	}

	// Check rate limiting settings
	if cfg.RateLimiting.Enabled != false {
		t.Errorf("Expected RateLimiting.Enabled to be false, got '%v'", cfg.RateLimiting.Enabled)
	}
	if cfg.RateLimiting.RequestsPerSecond != 100 {
		t.Errorf("Expected RateLimiting.RequestsPerSecond to be 100, got '%d'", cfg.RateLimiting.RequestsPerSecond)
	}
	if cfg.RateLimiting.Burst != 50 {
		t.Errorf("Expected RateLimiting.Burst to be 50, got '%d'", cfg.RateLimiting.Burst)
	}

	// Check caching settings
	if cfg.Caching.Enabled != false {
		t.Errorf("Expected Caching.Enabled to be false, got '%v'", cfg.Caching.Enabled)
	}
	if cfg.GetCachingDefaultTTL() != 300*time.Second {
		t.Errorf("Expected Caching.DefaultTTL to be 300s, got '%s'", cfg.GetCachingDefaultTTL())
	}
	if cfg.Caching.MaxSizeMB != 100 {
		t.Errorf("Expected Caching.MaxSizeMB to be 100, got '%d'", cfg.Caching.MaxSizeMB)
	}

	// Test JSON log format
	var buf bytes.Buffer
	handler := cfg.GetLogFormat(&buf)
	if _, ok := handler.(*slog.JSONHandler); !ok {
		t.Errorf("Expected JSONHandler, got %T", handler)
	}

	// Test default (text) log format
	cfg.Logging.Format = "text"
	handler = cfg.GetLogFormat(&buf)
	if _, ok := handler.(*slog.TextHandler); !ok {
		t.Errorf("Expected TextHandler, got %T", handler)
	}
}
