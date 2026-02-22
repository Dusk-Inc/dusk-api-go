package contracts

import "context"

type MetricsCollector interface {
	Collect(context context.Context) (contentType string, payload []byte, collectError error)
}

type MetricsRouterConfig struct {
	Collector MetricsCollector
}
