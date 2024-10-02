package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/shammianand/goproxy/internal/config"
	"github.com/shammianand/goproxy/internal/proxy"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Failed to run server", "error", err)
		os.Exit(1)
	}
}

func run() error {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	logger := slog.New(cfg.GetLogFormat())
	slog.SetDefault(logger)

	proxy, err := proxy.NewProxy(cfg.Proxy.TargetAddr, logger)
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:         cfg.Server.ListenAddr,
		Handler:      proxy,
		ReadTimeout:  cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:  cfg.Server.IdleTimeout * time.Second,
	}

	slog.Info("Starting GoProxy",
		"listen_addr", cfg.Server.ListenAddr,
		"target_addr", cfg.Proxy.TargetAddr,
		"log_level", cfg.Logging.Level,
	)

	return server.ListenAndServe()
}
