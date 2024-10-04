package integration

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

func TestMainIntegration(t *testing.T) {
	testCases := []struct {
		name                 string
		loadBalancingEnabled bool
	}{
		{"WithoutLoadBalancing", false},
		{"WithLoadBalancing", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := &config.Config{}
			cfg.Logging.Level = "debug"
			cfg.Logging.Format = "json"

			var logBuf bytes.Buffer
			logger := logger.New(cfg)
			jsonHandler := slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{Level: slog.LevelDebug})
			logger.Logger = slog.New(jsonHandler)

			backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello from backend"))
			}))
			defer backend.Close()

			var lb loadbalancer.LoadBalancer
			var err error

			if tc.loadBalancingEnabled {
				cfg.LoadBalancing.Enabled = true
				cfg.LoadBalancing.Algorithm = "round_robin"
				cfg.LoadBalancing.Backends = []string{backend.URL}

				lb, err = cfg.CreateLoadBalancer()
				if err != nil {
					t.Fatalf("Failed to create load balancer: %v", err)
				}
			}

			var proxyHandler http.Handler
			if tc.loadBalancingEnabled {
				proxyHandler, err = proxy.NewProxy("", lb, logger)
			} else {
				proxyHandler, err = proxy.NewProxy(backend.URL, nil, logger)
			}
			if err != nil {
				t.Fatalf("Failed to create proxy: %v", err)
			}

			server := httptest.NewServer(proxyHandler)
			defer server.Close()

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

			logOutput := logBuf.String()
			t.Logf("Log output: %s", logOutput)

			expectedLogs := []string{
				"Incoming request",
				"Response received",
			}

			for _, expected := range expectedLogs {
				if !strings.Contains(logOutput, expected) {
					t.Errorf("Expected log output to contain '%s', but it didn't", expected)
				}
			}
		})
	}
}
