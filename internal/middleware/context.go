package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"rps/internal/logging"

	"github.com/google/uuid"
)

type contextKey string

const (
	RequestIDKey contextKey = "requestID"
)

func GetRequestID(ctx context.Context) string {
	id, _ := ctx.Value(RequestIDKey).(string)
	return id
}

func Context(baseLogger *slog.Logger) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			requestID := uuid.NewString()

			// add requestID to context
			ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

			// add logger context
			logger := logging.WithTraceContext(ctx, baseLogger).
				With(logging.FieldRequestID, requestID)

			// add request_id to logger
			// bind logger to current context (Context-Aware Logging)
			ctx = logging.IntoContext(ctx, logger)

			// add request id to response header
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))

		})
	}
}
