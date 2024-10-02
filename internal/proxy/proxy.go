package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
	logger *slog.Logger
}

func NewProxy(target string, logger *slog.Logger) (*Proxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Customize the ReverseProxy director
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		logRequest(logger, req)
	}

	// Add a custom transport
	proxy.Transport = &loggingRoundTripper{
		logger: logger,
		next:   http.DefaultTransport,
	}

	return &Proxy{
		target: targetURL,
		proxy:  proxy,
		logger: logger,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}

func logRequest(logger *slog.Logger, r *http.Request) {
	logger.Info("Incoming request",
		"method", r.Method,
		"url", r.URL.String(),
		"remote_addr", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	)
}

type loggingRoundTripper struct {
	logger *slog.Logger
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
