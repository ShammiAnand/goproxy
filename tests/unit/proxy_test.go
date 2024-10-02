package unit

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shammianand/goproxy/internal/config"
	"github.com/shammianand/goproxy/internal/proxy"
	"github.com/shammianand/goproxy/pkg/logger"
)

func TestProxyLogging(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{}
	cfg.Logging.Level = "debug"
	cfg.Logging.Format = "json"

	// Create a buffer to capture log output
	var logBuffer bytes.Buffer
	logger := logger.New(cfg)

	// Create a JSON handler writing to the buffer
	jsonHandler := slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger.Logger = slog.New(jsonHandler)

	// Create a test server to act as the backend
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, client"))
	}))
	defer backend.Close()

	// Create our proxy
	proxy, err := proxy.NewProxy(backend.URL, logger)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	req.Header.Set("User-Agent", "test-agent")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Serve the request through our proxy
	proxy.ServeHTTP(rr, req)

	// Check the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := "Hello, client"
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Check the logs
	logOutput := logBuffer.String()

	expectedLogs := []string{
		"Incoming request",
		"method",
		"url",
		"remote_addr",
		"user_agent",
		"Response received",
		"status",
		"duration_ms",
		"content_length",
	}

	for _, expected := range expectedLogs {
		if !strings.Contains(logOutput, expected) {
			t.Errorf("Log output doesn't contain expected string: %s", expected)
		}
	}
}
