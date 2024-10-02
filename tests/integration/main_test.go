package integration

import (
	"bytes"
	"github.com/shammianand/goproxy/internal/config"
	"github.com/shammianand/goproxy/internal/proxy"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMainIntegration(t *testing.T) {
	// Create a temporary config file
	configContent := []byte(`
    server:
      listen_addr: ":8080"
      read_timeout: 5
      write_timeout: 10
      idle_timeout: 120
    proxy:
      target_addr: "http://localhost:8000"
      max_idle_conns: 100
      dial_timeout: 10
    logging:
      level: "info"
      format: "json"
  `)
	tmpfile, err := os.CreateTemp("", "config*.yaml")
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

	// Load the config
	cfg, err := config.Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create a buffer to capture log output
	var logBuf bytes.Buffer
	logger := slog.New(cfg.GetLogFormat(&logBuf))

	// Create a test backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello from backend"))
	}))
	defer backend.Close()

	// Override the target address in the config
	cfg.Proxy.TargetAddr = backend.URL

	// Create and start the proxy
	proxy, err := proxy.NewProxy(cfg.Proxy.TargetAddr, logger)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// Create a test server using our proxy
	server := httptest.NewServer(proxy)
	defer server.Close()

	// Make a test request
	resp, err := http.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if string(body) != "Hello from backend" {
		t.Errorf("Unexpected response body: %s", string(body))
	}

	// Check log output
	logOutput := logBuf.String()
	expectedLogs := []string{
		"Incoming request",
		"Response received",
	}

	for _, expected := range expectedLogs {
		if !bytes.Contains([]byte(logOutput), []byte(expected)) {
			t.Errorf("Expected log output to contain '%s', but it didn't", expected)
		}
	}
}
