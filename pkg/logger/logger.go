package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/shammianand/goproxy/internal/config"
)

// Logger wraps slog.Logger to provide additional methods if needed
type Logger struct {
	*slog.Logger
}

// New creates a new Logger instance
func New(cfg *config.Config) *Logger {
	var w io.Writer = os.Stdout
	// TODO: add file logging here if needed

	handler := cfg.GetLogFormat(w)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}

// Named returns a new Logger with the given name added to the logger's context
func (l *Logger) Named(name string) *Logger {
	return &Logger{Logger: l.Logger.With("logger", name)}
}

// With returns a new Logger with the given key-value pairs added to the logger's context
func (l *Logger) With(args ...any) *Logger {
	return &Logger{Logger: l.Logger.With(args...)}
}

// Debug logs at Debug level
func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

// Info logs at Info level
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Warn logs at Warn level
func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

// Error logs at Error level
func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}
