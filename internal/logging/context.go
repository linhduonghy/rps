package logging

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

const (
	FieldRequestID = "request_id"
	FieldTraceID   = "trace_id"
	FieldSpanID    = "span_id"

	FieldMethod = "method"
	FieldPath   = "path"
)

type contextKey struct{}

func IntoContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, logger)
}

func WithTraceContext(ctx context.Context, logger *slog.Logger) *slog.Logger {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()

	if !sc.IsValid() {
		return logger
	}

	return logger.With(
		FieldTraceID, sc.TraceID().String(),
		FieldSpanID, sc.SpanID().String(),
	)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(contextKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}
