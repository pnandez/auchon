// Package logger builds the kernel's structured logger, the only way to log.
package logger

import (
	"io"
	"log/slog"
)

// New returns a JSON logger writing to w at the given level.
func New(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level}))
}
