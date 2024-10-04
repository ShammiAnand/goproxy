package unit

import (
	"net/url"
	"testing"

	"github.com/shammianand/goproxy/internal/loadbalancer"
)

func TestRoundRobinBalancer(t *testing.T) {

	backends := []*loadbalancer.Backend{
		{URL: mustParseURL("http://backend1.com"), Healthy: true},
		{URL: mustParseURL("http://backend2.com"), Healthy: true},
		{URL: mustParseURL("http://backend3.com"), Healthy: true},
	}
	balancer := loadbalancer.NewRoundRobinBalancer(backends)

	// Test round-robin selection
	seenBackends := make(map[string]bool)
	for i := 0; i < 9; i++ {
		backend, err := balancer.NextBackend()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		seenBackends[backend.URL.String()] = true
	}
	if len(seenBackends) != 3 {
		t.Errorf("Expected to see 3 unique backends, got %d", len(seenBackends))
	}

	// Test unhealthy backend
	balancer.HealthCheck(backends[1], false)
	for i := 0; i < 4; i++ {
		backend, err := balancer.NextBackend()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if backend.URL == backends[1].URL {
			t.Errorf("Selected unhealthy backend: %s", backend.URL)
		}
	}

	// Test no healthy backends
	balancer.HealthCheck(backends[0], false)
	balancer.HealthCheck(backends[2], false)
	_, err := balancer.NextBackend()
	if err != loadbalancer.ErrNoHealthyBackends {
		t.Errorf("Expected ErrNoHealthyBackends, got %v", err)
	}

	// Test update backends
	newBackends := []*loadbalancer.Backend{
		{URL: mustParseURL("http://newbackend1.com"), Healthy: true},
		{URL: mustParseURL("http://newbackend2.com"), Healthy: true},
	}
	balancer.UpdateBackends(newBackends)
	if len(balancer.Backends()) != 2 {
		t.Errorf("Expected 2 backends, got %d", len(balancer.Backends()))
	}
}

func mustParseURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}
	return u
}
