package performance

import (
	"github.com/shammianand/goproxy/internal/proxy"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkProxy(b *testing.B) {
	logger := slog.New(slog.NewJSONHandler(nil, nil))
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := proxy.NewProxy(backend.URL, logger)
	if err != nil {
		b.Fatalf("Failed to create proxy: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "http://example.com/test", nil)
		w := httptest.NewRecorder()
		proxy.ServeHTTP(w, req)
	}
}
