package config

import (
	"log/slog"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		ListenAddr   string        `yaml:"listen_addr"`
		ReadTimeout  time.Duration `yaml:"read_timeout"`
		WriteTimeout time.Duration `yaml:"write_timeout"`
		IdleTimeout  time.Duration `yaml:"idle_timeout"`
	} `yaml:"server"`
	Proxy struct {
		TargetAddr   string        `yaml:"target_addr"`
		MaxIdleConns int           `yaml:"max_idle_conns"`
		DialTimeout  time.Duration `yaml:"dial_timeout"`
	} `yaml:"proxy"`
	LoadBalancing struct {
		Enabled   bool     `yaml:"enabled"`
		Algorithm string   `yaml:"algorithm"`
		Backends  []string `yaml:"backends"`
	} `yaml:"load_balancing"`
	TLS struct {
		Enabled  bool   `yaml:"enabled"`
		CertFile string `yaml:"cert_file"`
		KeyFile  string `yaml:"key_file"`
	} `yaml:"tls"`
	Logging struct {
		Level  string `yaml:"level"`
		Format string `yaml:"format"`
	} `yaml:"logging"`
	Metrics struct {
		Enabled bool   `yaml:"enabled"`
		Address string `yaml:"address"`
	} `yaml:"metrics"`
	RateLimiting struct {
		Enabled           bool `yaml:"enabled"`
		RequestsPerSecond int  `yaml:"requests_per_second"`
		Burst             int  `yaml:"burst"`
	} `yaml:"rate_limiting"`
	Caching struct {
		Enabled    bool          `yaml:"enabled"`
		DefaultTTL time.Duration `yaml:"default_ttl"`
		MaxSizeMB  int           `yaml:"max_size_mb"`
	} `yaml:"caching"`
}

func Load(configPath string) (*Config, error) {
	config := &Config{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) GetLogLevel() slog.Level {
	switch c.Logging.Level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (c *Config) GetLogFormat() slog.Handler {
	switch c.Logging.Format {
	case "json":
		return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: c.GetLogLevel()})
	default:
		return slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: c.GetLogLevel()})
	}
}
