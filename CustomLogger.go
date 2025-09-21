package mlogger

import (
	"context"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

// multiHandler is a custom slog.Handler that forwards log records to multiple handlers.
type multiHandler struct {
	handlers []slog.Handler
}

// NewMultiHandler creates a new multiHandler with the given handlers.
func NewMultiHandler(handlers ...slog.Handler) *multiHandler {
	return &multiHandler{handlers: handlers}
}

// Handle forwards the log record to all contained handlers.
func (h *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// Enabled returns true if any of the contained handlers are enabled for the given level.
func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// WithAttrs returns a new multiHandler with attributes added to all contained handlers.
func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

// WithGroup returns a new multiHandler with a group added to all contained handlers.
func (h *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		newHandlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

func InitMultiLogger() {
	// Configure the lumberjack logger
	rotator := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    1,    // rotate after 10 megabytes
		MaxBackups: 5,    // keep at most 5 old log files
		MaxAge:     30,   // keep log files for up to 30 days
		Compress:   true, // compress rotated log files
	}

	// Create a JSON handler that writes to our lumberjack rotator
	fileHandler := slog.NewJSONHandler(rotator, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// 2. Create a handler to write to the console (Text format)
	consoleHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	// 3. Combine the handlers using our custom multiHandler
	multi := NewMultiHandler(consoleHandler, fileHandler)
	logger := slog.New(multi)

	// 4. Set the new logger as the default
	slog.SetDefault(logger)
}
