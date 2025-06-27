package logutil

import (
	"context"
	"log/slog"
)

type ctxKey string

const traceIDKey ctxKey = "trace_id"

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func TraceIDFromContext(ctx context.Context) string {
	if val := ctx.Value(traceIDKey); val != nil {
		if id, ok := val.(string); ok {
			return id
		}
	}
	return "unknown"
}

func Logger(ctx context.Context) *slog.Logger {
	return slog.Default().With("trace_id", TraceIDFromContext(ctx))
}
