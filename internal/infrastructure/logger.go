package infrastructure

import (
	"context"
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger
}

// New creates a new logger instance with the specified log level
func New(level slog.Level) *Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &Logger{
		Logger: logger,
	}
}

// WithContext adds context fields to the logger
func (l *Logger) WithContext(ctx context.Context) *slog.Logger {
	// Add any context-specific fields here
	// For example, request ID, user ID, etc.
	return l.Logger
}

// With adds key-value pairs to the logger
func (l *Logger) With(args ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// Debug logs a debug message with optional key-value pairs
func (l *Logger) Debug(msg string, args ...any) {
	l.Logger.Debug(msg, args...)
}

// Info logs an info message with optional key-value pairs
func (l *Logger) Info(msg string, args ...any) {
	l.Logger.Info(msg, args...)
}

// Warn logs a warning message with optional key-value pairs
func (l *Logger) Warn(msg string, args ...any) {
	l.Logger.Warn(msg, args...)
}

// Error logs an error message with optional key-value pairs
func (l *Logger) Error(msg string, args ...any) {
	l.Logger.Error(msg, args...)
}

// Fatal logs a fatal message and exits the application
func (l *Logger) Fatal(msg string, args ...any) {
	l.Logger.Error(msg, args...)
	os.Exit(1)
}
