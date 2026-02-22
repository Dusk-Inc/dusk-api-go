package contracts

import "context"

type ReadinessCheck func(context context.Context) bool

type HealthRouterConfig struct {
	Readiness ReadinessCheck
}
