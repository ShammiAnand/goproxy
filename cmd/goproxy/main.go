package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"

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

	return RunWithConfig(*configPath)
}

// RunWithConfig is exported for testing purposes
func RunWithConfig(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.GetLogLevel(),
	}))
	slog.SetDefault(logger)

	proxy, err := proxy.NewProxy(cfg.TargetAddr, logger)
	if err != nil {
		return err
	}

	slog.Info("Starting GoProxy", "listen_addr", cfg.ListenAddr)
	return http.ListenAndServe(cfg.ListenAddr, proxy)
}
