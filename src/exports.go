package src

import (
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/functions"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/modules"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/routes"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/tokens"
)

type (
	AppManager             = modules.AppManager
	HealthRouter           = routes.HealthRouter
	MetricsRouter          = routes.MetricsRouter
	RuntimeManager         = modules.RuntimeManager
	SecretManager          = modules.SecretManager
	WellKnownRouter        = modules.WellKnownRouter
	ServiceDecorator       = modules.ServiceDecorator
	ServiceDecoratorConfig = contracts.ServiceDecoratorConfig
)

var (
	ServiceDecoratorPhase = tokens.ServiceDecoratorPhase
)

func NewAppManager(config contracts.AppManagerConfig) *modules.AppManager {
	return modules.NewAppManager(config)
}

func NewHealthRouter(config contracts.HealthRouterConfig) *routes.HealthRouter {
	return routes.NewHealthRouter(config)
}

func NewMetricsRouter(config contracts.MetricsRouterConfig) *routes.MetricsRouter {
	return routes.NewMetricsRouter(config)
}

func NewWellKnownRouter(config contracts.WellKnownRouterConfig) *modules.WellKnownRouter {
	return modules.NewWellKnownRouter(config)
}

func ParseEnv() (functions.EnvModel, error) {
	return functions.ParseEnv()
}
