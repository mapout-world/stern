package handler

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/mapout-world/stern/lib/metrics"
)

type Metrics struct {
	ExampleCounter prometheus.Counter
}

func (h *Handler) RegisterMetrics() {
	h.metrics.ExampleCounter = metrics.NewCounter("example", metrics.WithHelp("Example metric."))
}
