package infrastructure

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
)

type CustomSourceHandler struct {
	Handler slog.Handler
}

func (h *CustomSourceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *CustomSourceHandler) Handle(ctx context.Context, r slog.Record) error {
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		// Format: cmd/main.go:42
		source := filepath.Base(filepath.Dir(f.File)) + "/" + filepath.Base(f.File) + ":" + fmt.Sprint(f.Line)
		r.AddAttrs(slog.String("source", source))
	}
	return h.Handler.Handle(ctx, r)
}

func (h *CustomSourceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomSourceHandler{h.Handler.WithAttrs(attrs)}
}

func (h *CustomSourceHandler) WithGroup(name string) slog.Handler {
	return &CustomSourceHandler{h.Handler.WithGroup(name)}
}
