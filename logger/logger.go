package logger

import (
	"context"
	"log/slog"
)

type TraceIDKey struct{}

func ContextLogger(ctx context.Context) *slog.Logger {
	traceID := ctx.Value(TraceIDKey{})
	if traceID == nil {
		return slog.Default()
	}
	return slog.With("trace_id", traceID.(string))
}
