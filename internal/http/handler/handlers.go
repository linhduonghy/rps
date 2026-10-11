package handler

import "rps/internal/metric"

type Handlers struct {
	Metric *metric.HttpHandler
	Policy *PolicyHandler
}
