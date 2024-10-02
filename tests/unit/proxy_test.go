package unit

import (
	"github.com/shammianand/goproxy/internal/proxy"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProxyServeHTTP(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(nil, nil))
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	proxy, err := proxy.NewProxy(backend.URL, logger)
	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()

	proxy.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
