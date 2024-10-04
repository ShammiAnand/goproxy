package simulation

import (
	// "fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/shammianand/goproxy/internal/config"
	"github.com/shammianand/goproxy/internal/loadbalancer"
	p "github.com/shammianand/goproxy/internal/proxy"
	"github.com/shammianand/goproxy/pkg/logger"
	"golang.org/x/exp/rand"
)

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}

func TestLoadBalancerSimulation(t *testing.T) {
	cfg := &config.Config{}
	cfg.Logging.Level = "debug" // Change to debug for more information
	cfg.Logging.Format = "json"

	log := logger.New(cfg)

	t.Run("WithLoadBalancer", func(t *testing.T) {
		runLoadBalancerTest(t, log, true)
	})

	t.Run("WithoutLoadBalancer", func(t *testing.T) {
		runLoadBalancerTest(t, log, false)
	})
}

func runLoadBalancerTest(t *testing.T, log *logger.Logger, useLoadBalancer bool) {
	// Create mock backends
	backend1 := createMockBackend("Backend 1")
	defer backend1.Close()
	backend2 := createMockBackend("Backend 2")
	defer backend2.Close()
	backend3 := createMockBackend("Backend 3")
	defer backend3.Close()

	var proxy http.Handler
	var err error
	var balancer *loadbalancer.RoundRobinBalancer

	if useLoadBalancer {
		// Create load balancer
		backends := []*loadbalancer.Backend{
			{URL: mustParseURL(backend1.URL), Healthy: true},
			{URL: mustParseURL(backend2.URL), Healthy: true},
			{URL: mustParseURL(backend3.URL), Healthy: true},
		}
		balancer = loadbalancer.NewRoundRobinBalancer(backends)

		// Print initial state of load balancer
		t.Logf("Initial load balancer state: %+v", balancer)

		// Create proxy with load balancer
		proxy, err = p.NewProxy("", balancer, log)
	} else {
		// Create proxy without load balancer
		proxy, err = p.NewProxy(backend1.URL, nil, log)
	}

	if err != nil {
		t.Fatalf("Failed to create proxy: %v", err)
	}

	// Create test server using our proxy
	server := httptest.NewServer(proxy)
	defer server.Close()

	// Simulate traffic
	requestCount := 1000
	responseCounts := simulateTraffic(t, server, requestCount)

	// Check distribution
	checkDistribution(t, responseCounts, requestCount, useLoadBalancer)

	if useLoadBalancer {
		// Print load balancer state after traffic
		t.Logf("Load balancer state after traffic: %+v", balancer)

		// Simulate a backend going down
		balancer.HealthCheck(balancer.Backends()[1], false)
		responseCounts = simulateTraffic(t, server, requestCount)

		// Check distribution with one backend down
		checkDistributionWithBackendDown(t, responseCounts, requestCount)
	}
}

func simulateTraffic(t *testing.T, server *httptest.Server, requestCount int) map[string]int {
	var wg sync.WaitGroup
	responseCounts := make(map[string]int)
	var mu sync.Mutex

	for i := 0; i < requestCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resp, err := http.Get(server.URL)
			if err != nil {
				t.Errorf("Request failed: %v", err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %v", err)
				return
			}
			mu.Lock()
			responseCounts[string(body)]++
			mu.Unlock()
		}()
		// Add some randomness to request timing
		time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	}
	wg.Wait()
	return responseCounts
}

func checkDistribution(t *testing.T, responseCounts map[string]int, requestCount int, useLoadBalancer bool) {
	t.Logf("Total requests: %d", requestCount)
	for backend, count := range responseCounts {
		percentage := float64(count) / float64(requestCount) * 100
		t.Logf("%s: %d requests (%.2f%%)", backend, count, percentage)
		if useLoadBalancer {
			if count == 0 {
				t.Errorf("Backend %s received no requests", backend)
			}
		} else {
			if count != requestCount {
				t.Errorf("Expected all requests to go to single backend, got %d out of %d", count, requestCount)
			}
		}
	}
	if useLoadBalancer && len(responseCounts) != 3 {
		t.Errorf("Expected responses from 3 backends, got %d", len(responseCounts))
	}
}

func checkDistributionWithBackendDown(t *testing.T, responseCounts map[string]int, requestCount int) {
	t.Logf("Total requests with one backend down: %d", requestCount)
	for backend, count := range responseCounts {
		percentage := float64(count) / float64(requestCount) * 100
		t.Logf("%s: %d requests (%.2f%%)", backend, count, percentage)
		if backend == "Backend 2" && count > 0 {
			t.Errorf("Received requests for unhealthy backend: %s", backend)
		}
		if backend == "Backend 1" || backend == "Backend 3" {
			if count == 0 {
				t.Errorf("Backend %s received no requests", backend)
			}
		}
	}
	if len(responseCounts) != 2 {
		t.Errorf("Expected responses from 2 backends, got %d", len(responseCounts))
	}
}

func createMockBackend(name string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(name))
	}))
}
