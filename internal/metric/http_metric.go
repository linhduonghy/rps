package metric

import (
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func newHttpRequestStatusCounterMetric() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "app_http_requests_total", Help: "Total HTTP requests"},
		[]string{"method", "path", "status"},
	)
}

func newHttpDurationMetric() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_http_duration_seconds",
			Help:    "HTTP request duration.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{
			"method",
			"path",
		},
	)
}

func (m *Metrics) IncHttpRequestStatusCounter(method string, path string, status int) {
	m.httpRequests.WithLabelValues(
		method,
		path,
		strconv.Itoa(status),
	).Inc()
}

func (m *Metrics) SetHttpDuration(method string, path string, duration float64) {
	m.httpDuration.WithLabelValues(
		method, path).Observe(duration)
}
