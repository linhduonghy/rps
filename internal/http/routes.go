package http

import (
	"log/slog"
	"net/http"
	"rps/internal/http/handler"
	"rps/internal/metric"
	"rps/internal/middleware"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewRouter(
	logger *slog.Logger,
	h *handler.Handlers,
	metrics *metric.Metrics,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recover(logger))
	r.Use(middleware.Context(logger))
	r.Use(middleware.Metric(metrics))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// metrics
	r.Handle("/metrics", h.Metric)

	// app route
	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/policies", func(r chi.Router) {
			r.Get("/", h.Policy.ListPolicies)
			r.Get("/{resourceType}/{resourceId}", h.Policy.GetPolicy)
			r.Post("/", h.Policy.CreatePolicy)
			r.Put("/{resourceType}/{resourceId}", h.Policy.UpdatePolicy)
			r.Delete("/{resourceType}/{resourceId}", h.Policy.DeletePolicy)
		})
	})
	// add otelhttp
	return otelhttp.NewHandler(
		r,
		"http.server",
	)
}
