package middleware

import (
	"net/http"
	"rps/internal/logging"
	"rps/internal/metric"
	"time"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	status int
}

func NewStatusResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (sw *CustomResponseWriter) WriteHeader(status int) {
	sw.status = status
	sw.ResponseWriter.WriteHeader(status)
}

func Metric(m *metric.Metrics) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := NewStatusResponseWriter(w)
			next.ServeHTTP(sw, r)

			duration := time.Since(start)

			routePattern := r.Pattern
			if routePattern == "" {
				routePattern = "unknown"
			}

			logging.FromContext(r.Context()).Info(
				"http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", sw.status,
				"duration", duration,
			)

			m.IncHttpRequestStatusCounter(r.Method, routePattern, sw.status)
			m.SetHttpDuration(r.Method, routePattern, duration.Seconds())

		})
	}
}
