package contracts

import "log/slog"

type AppManagerConfig struct {
	ServiceName string
	Logger      *slog.Logger
	Readiness   ReadinessCheck
	// Collector, if non-nil, replaces the default runtime metrics collector
	// registered at /metrics. Implementations emit Prometheus text format.
	Collector MetricsCollector
}
