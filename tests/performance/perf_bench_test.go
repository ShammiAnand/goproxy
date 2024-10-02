package performance

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shammianand/goproxy/internal/config"
	"github.com/shammianand/goproxy/internal/proxy"
	"github.com/shammianand/goproxy/pkg/logger"
)

func BenchmarkProxy(b *testing.B) {
	// Create a null logger to avoid logging overhead during benchmarking
	// Create a test configuration
	cfg := &config.Config{}
	cfg.Logging.Level = "debug"
	cfg.Logging.Format = "json"

	// Create a buffer to capture log output
	var logBuffer bytes.Buffer
	logger := logger.New(cfg)
	jsonHandler := slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger.Logger = slog.New(jsonHandler)

	// Create a test backend server
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, client"))
	}))
	defer backend.Close()

	// Create the proxy
	proxy, err := proxy.NewProxy(backend.URL, logger)
	if err != nil {
		b.Fatalf("Failed to create proxy: %v", err)
	}

	// Create a test request
	req := httptest.NewRequest("GET", "http://example.com/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		proxy.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			b.Fatalf("Unexpected status code: %d", w.Code)
		}
	}
}
