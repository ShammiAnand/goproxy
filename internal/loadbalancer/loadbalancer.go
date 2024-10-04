package loadbalancer

import (
	"errors"
	"net/url"
	// "sync"
)

// Backend represents a backend server
type Backend struct {
	URL     *url.URL
	Healthy bool
}

// LoadBalancer interface defines the methods a load balancer should implement
type LoadBalancer interface {
	NextBackend() (*Backend, error)
	UpdateBackends(backends []*Backend)
	HealthCheck(backend *Backend, healthy bool)
	Backends() []*Backend
}

var (
	ErrNoHealthyBackends = errors.New("no healthy backends available")
)
