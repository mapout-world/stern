package metrics

import (
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/mapout-world/stern/lib/config"
)

var normalizeNamespaceRe = regexp.MustCompile(`[^A-z0-9_]`)

func serviceNamespace() string {
	service := config.Get("name").String("unknown")
	return "mapout_" + normalizeNamespaceRe.ReplaceAllString(service, "_")
}

func NewCounter(name string, opts ...MetricOption) prometheus.Counter {
	options := prometheus.Opts{
		Namespace: serviceNamespace(),
		Name:      name,
	}

	for _, opt := range opts {
		opt(&options)
	}

	return promauto.NewCounter(prometheus.CounterOpts(options))
}

type MetricOption func(*prometheus.Opts)

func WithLabels(labels prometheus.Labels) MetricOption {
	return func(o *prometheus.Opts) {
		o.ConstLabels = labels
	}
}

func WithHelp(help string) MetricOption {
	return func(o *prometheus.Opts) {
		o.Help = help
	}
}
