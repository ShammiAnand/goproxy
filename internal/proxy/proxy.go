package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/shammianand/goproxy/internal/loadbalancer"
	"github.com/shammianand/goproxy/pkg/logger"
)

type Proxy struct {
	target       *url.URL
	proxy        *httputil.ReverseProxy
	logger       *logger.Logger
	loadBalancer loadbalancer.LoadBalancer
}

func NewProxy(target string, lb loadbalancer.LoadBalancer, logger *logger.Logger) (http.Handler, error) {
	var targetURL *url.URL
	var err error
	if target != "" {
		targetURL, err = url.Parse(target)
		if err != nil {
			return nil, err
		}
	}

	p := &Proxy{
		target:       targetURL,
		loadBalancer: lb,
		logger:       logger,
	}

	if lb == nil && targetURL != nil {
		p.proxy = httputil.NewSingleHostReverseProxy(targetURL)
		p.proxy.ErrorLog = slog.NewLogLogger(logger.Handler(), slog.LevelError)
	}

	return p, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var proxyToUse *httputil.ReverseProxy
	var backendURL *url.URL

	if p.loadBalancer != nil {
		backend, err := p.loadBalancer.NextBackend()
		if err != nil {
			p.logger.Error("Failed to get next backend", "error", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		backendURL = backend.URL
		proxyToUse = httputil.NewSingleHostReverseProxy(backendURL)
	} else if p.proxy != nil {
		proxyToUse = p.proxy
		backendURL = p.target
	} else {
		p.logger.Error("No backend or load balancer configured")
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
		return
	}

	proxyToUse.ErrorLog = slog.NewLogLogger(p.logger.Handler(), slog.LevelError)

	// Modify the request to match the backend URL
	r.URL.Host = backendURL.Host
	r.URL.Scheme = backendURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = backendURL.Host

	// Log the incoming request
	p.logger.Info("Incoming request",
		"method", r.Method,
		"url", r.URL.String(),
		"backend", backendURL.String(),
	)

	// Capture the response
	rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
	proxyToUse.ServeHTTP(rw, r)

	// Log the response
	p.logger.Info("Response received",
		"status", rw.statusCode,
		"backend", backendURL.String(),
	)
}

// responseWriter is a custom ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func logRequest(logger *logger.Logger, r *http.Request) {
	logger.Info("Incoming request",
		"method", r.Method,
		"url", r.URL.String(),
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	)
}

type loggingRoundTripper struct {
	logger *logger.Logger
	next   http.RoundTripper
}

func (l *loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := l.next.RoundTrip(req)
	duration := time.Since(start)

	if err != nil {
		l.logger.Error("Error in round trip",
			"error", err,
			"duration_ms", duration.Milliseconds(),
		)
		return nil, err
	}

	l.logger.Info("Response received",
		"status", resp.Status,
		"duration_ms", duration.Milliseconds(),
		"content_length", resp.ContentLength,
	)

	return resp, nil
}
