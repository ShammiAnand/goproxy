package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	ListenAddr string `yaml:"listen_addr"`
	TargetAddr string `yaml:"target_addr"`
	LogLevel   string `yaml:"log_level"`
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
	switch c.LogLevel {
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
