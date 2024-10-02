package integration

import (
	"bytes"
	"github.com/shammianand/goproxy/internal/config"
	"github.com/shammianand/goproxy/internal/proxy"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	// Setup a mock config
	configContent := []byte(`
    listen_addr: ":8080"
    target_addr: "http://example.com"
    log_level: "info"
  `)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp config file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(configContent); err != nil {
		t.Fatalf("Failed to write to temp config file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Failed to close temp config file: %v", err)
	}

	// Capture log output
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, nil))
	slog.SetDefault(logger)

	// Create a test server to simulate the target
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()

	// Override the target address in the config
	cfg, err := config.Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	cfg.TargetAddr = testServer.URL

	// Create and start the proxy
	p, err := proxy.NewProxy(cfg.TargetAddr, logger)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	proxyServer := httptest.NewServer(p)
	defer proxyServer.Close()

	// Make a test request
	resp, err := http.Get(proxyServer.URL + "/test")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Check log output
	logOutput := logBuf.String()
	if !bytes.Contains([]byte(logOutput), []byte("Received request")) {
		t.Errorf("Expected log output to contain 'Received request', got %s", logOutput)
	}
}
