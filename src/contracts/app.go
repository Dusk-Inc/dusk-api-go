package contracts

import "log/slog"

type AppManagerConfig struct {
	ServiceName string
	Logger      *slog.Logger
	Readiness   ReadinessCheck
}
