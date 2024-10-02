package performance

import (
	"github.com/shammianand/goproxy/internal/proxy"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkProxy(b *testing.B) {
	// Create a null logger to avoid logging overhead during benchmarking
	logger := slog.New(slog.NewJSONHandler(nil, &slog.HandlerOptions{Level: slog.LevelError}))

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
