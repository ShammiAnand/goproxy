package performance

import (
	"bytes"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/shammianand/goproxy/internal/config"
	"github.com/shammianand/goproxy/internal/loadbalancer"
	"github.com/shammianand/goproxy/internal/proxy"
	"github.com/shammianand/goproxy/pkg/logger"
)

func BenchmarkProxy(b *testing.B) {
	// Run the benchmark with different concurrency levels
	for _, concurrency := range []int{1, 10, 50, 100} {
		b.Run(fmt.Sprintf("Concurrency-%d", concurrency), func(b *testing.B) {
			runBenchmark(b, concurrency)
		})
	}
}

func runBenchmark(b *testing.B, concurrency int) {
	b.Logf("Starting benchmark with %d iterations and concurrency %d", b.N, concurrency)

	cfg := &config.Config{}
	cfg.Logging.Level = "error"
	cfg.Logging.Format = "json"

	var logBuffer bytes.Buffer
	logger := logger.New(cfg)
	jsonHandler := slog.NewJSONHandler(&logBuffer, &slog.HandlerOptions{Level: slog.LevelError})
	logger.Logger = slog.New(jsonHandler)

	backends := createBackends(5) // Increase number of backends
	defer func() {
		for _, backend := range backends {
			backend.Close()
		}
	}()

	lbBackends := make([]*loadbalancer.Backend, len(backends))
	for i, backend := range backends {
		lbBackends[i] = &loadbalancer.Backend{URL: mustParseURL(backend.URL), Healthy: true}
	}
	lb := loadbalancer.NewRoundRobinBalancer(lbBackends)

	proxyHandler, err := proxy.NewProxy("", lb, logger)
	if err != nil {
		b.Fatalf("Failed to create proxy: %v", err)
	}

	b.ResetTimer()

	var wg sync.WaitGroup
	requestChan := make(chan int, b.N)
	resultChan := make(chan result, b.N)

	// Start worker goroutines
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(b, &wg, requestChan, resultChan, proxyHandler)
	}

	// Send requests
	for i := 0; i < b.N; i++ {
		requestChan <- i
	}
	close(requestChan)

	wg.Wait()
	close(resultChan)

	// Process results
	var totalLatency time.Duration
	var minLatency, maxLatency time.Duration
	statusCodes := make(map[int]int)
	var successCount int

	for res := range resultChan {
		totalLatency += res.latency
		statusCodes[res.statusCode]++

		if res.statusCode == http.StatusOK {
			successCount++
		}

		if res.latency < minLatency || minLatency == 0 {
			minLatency = res.latency
		}
		if res.latency > maxLatency {
			maxLatency = res.latency
		}
	}

	avgLatency := totalLatency / time.Duration(b.N)
	successRate := float64(successCount) / float64(b.N) * 100

	b.ReportMetric(float64(avgLatency.Nanoseconds())/1e6, "avg_latency_ms")
	b.ReportMetric(float64(minLatency.Nanoseconds())/1e6, "min_latency_ms")
	b.ReportMetric(float64(maxLatency.Nanoseconds())/1e6, "max_latency_ms")
	b.ReportMetric(successRate, "success_rate_%")

	b.Logf("Benchmark completed:")
	b.Logf("  Total requests: %d", b.N)
	b.Logf("  Successful requests: %d", successCount)
	b.Logf("  Success rate: %.2f%%", successRate)
	b.Logf("  Average latency: %v", avgLatency)
	b.Logf("  Min latency: %v", minLatency)
	b.Logf("  Max latency: %v", maxLatency)
	b.Logf("  Status code distribution:")
	for code, count := range statusCodes {
		b.Logf("    %d: %d (%.2f%%)", code, count, float64(count)/float64(b.N)*100)
	}
}

type result struct {
	statusCode int
	latency    time.Duration
}

func worker(b *testing.B, wg *sync.WaitGroup, requestChan <-chan int, resultChan chan<- result, proxyHandler http.Handler) {
	defer wg.Done()

	for range requestChan {
		req := httptest.NewRequest("GET", fmt.Sprintf("http://example.com/test?id=%d", rand.Intn(1000)), nil)
		w := httptest.NewRecorder()

		start := time.Now()
		proxyHandler.ServeHTTP(w, req)
		latency := time.Since(start)

		resultChan <- result{
			statusCode: w.Code,
			latency:    latency,
		}
	}
}

func createBackends(count int) []*httptest.Server {
	backends := make([]*httptest.Server, count)
	for i := 0; i < count; i++ {
		backends[i] = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate varying response times and occasional errors
			if rand.Float32() < 0.01 { // 1% chance of error
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(20)+1))
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Hello from backend %d", i+1)
		}))
	}
	return backends
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
