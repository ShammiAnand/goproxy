package unit

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shammianand/goproxy/internal/config"
	"github.com/shammianand/goproxy/internal/loadbalancer"
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

	// Create a load balancer with a single backend
	backends := []*loadbalancer.Backend{
		{URL: mustParseURL(backend.URL), Healthy: true},
	}
	loadBalancer := loadbalancer.NewRoundRobinBalancer(backends)

	// Create our proxy
	proxy, err := proxy.NewProxy(backend.URL, loadBalancer, logger)
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

	// Check the logs
	logOutput := logBuffer.String()

	expectedLogs := []string{
		"Incoming request",
		"method",
		"url",
		"backend",
	}

	for _, expected := range expectedLogs {
		if !strings.Contains(logOutput, expected) {
			t.Errorf("Log output doesn't contain expected string: %s", expected)
		}
	}

	// Check if the log output contains a valid URL
	if !strings.Contains(logOutput, "http://") {
		t.Errorf("Log output doesn't contain a valid URL")
	}
}

func TestProxyWithLoadBalancing(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{}
	cfg.Logging.Level = "debug"
	cfg.Logging.Format = "json"

	// Create mock backends
	backend1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from backend 1"))
	}))
	defer backend1.Close()

	backend2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Response from backend 2"))
	}))
	defer backend2.Close()

	// Create load balancer
	backends := []*loadbalancer.Backend{
		{URL: mustParseURL(backend1.URL), Healthy: true},
		{URL: mustParseURL(backend2.URL), Healthy: true},
	}
	balancer := loadbalancer.NewRoundRobinBalancer(backends)

	// Create logger
	var logBuffer bytes.Buffer
	log := logger.New(cfg)

	// Create a JSON handler writing to the buffer
	jsonHandler := slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug})
	log.Logger = slog.New(jsonHandler)

	// Create proxy
	proxy, err := proxy.NewProxy("", balancer, log)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// Create test server using our proxy
	server := httptest.NewServer(proxy)
	defer server.Close()

	// Make requests and check round-robin behavior
	responses := make(map[string]int)
	for i := 0; i < 4; i++ {
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		responses[string(body)]++
	}

	// Check that we got responses from both backends
	if len(responses) != 2 {
		t.Errorf("Expected responses from 2 backends, got %d", len(responses))
	}

	// Check that each backend was used twice
	for backend, count := range responses {
		if count != 2 {
			t.Errorf("Backend %s was used %d times, expected 2", backend, count)
		}
	}

	// Check logs
	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "Incoming request") {
		t.Error("Log output doesn't contain 'Incoming request'")
	}
	if !strings.Contains(logOutput, backend1.URL) || !strings.Contains(logOutput, backend2.URL) {
		t.Error("Log output doesn't contain backend URLs")
	}
}
