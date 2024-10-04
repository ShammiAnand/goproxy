package loadbalancer

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
)

// RoundRobinBalancer implements the LoadBalancer interface using a round-robin algorithm
type RoundRobinBalancer struct {
	backends []*Backend
	mutex    sync.RWMutex
	current  uint32
}

// NewRoundRobinBalancer creates a new RoundRobinBalancer
func NewRoundRobinBalancer(backends []*Backend) *RoundRobinBalancer {
	return &RoundRobinBalancer{
		backends: backends,
		current:  rand.Uint32(),
	}
}

// NextBackend returns the next backend using round-robin selection
func (r *RoundRobinBalancer) NextBackend() (*Backend, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if len(r.backends) == 0 {
		return nil, ErrNoHealthyBackends
	}

	next := int(atomic.AddUint32(&r.current, 1) % uint32(len(r.backends)))
	for i := 0; i < len(r.backends); i++ {
		idx := (next + i) % len(r.backends)
		if r.backends[idx].Healthy {
			log.Printf("Selected backend %d: %s", idx, r.backends[idx].URL)
			return r.backends[idx], nil
		}
	}

	return nil, ErrNoHealthyBackends
}

// UpdateBackends updates the list of backends
func (r *RoundRobinBalancer) UpdateBackends(backends []*Backend) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.backends = backends
}

// HealthCheck updates the health status of a backend
func (r *RoundRobinBalancer) HealthCheck(backend *Backend, healthy bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, b := range r.backends {
		if b.URL.String() == backend.URL.String() {
			b.Healthy = healthy
			break
		}
	}
}

// Backends returns the list of backends
func (r *RoundRobinBalancer) Backends() []*Backend {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.backends
}
