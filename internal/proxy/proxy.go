package proxy

import (
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
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

	return &Proxy{
		target: targetURL,
		proxy:  httputil.NewSingleHostReverseProxy(targetURL),
		logger: logger,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.logger.Info("Received request",
		"method", r.Method,
		"url", r.URL.String(),
		"remote_addr", r.RemoteAddr,
	)
	p.proxy.ServeHTTP(w, r)
}
