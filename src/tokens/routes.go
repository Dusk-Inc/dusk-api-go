package tokens

import "github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"

var HealthLiveContract = contracts.RouteContract{
	Method: contracts.RouteMethodGet,
	Path:   "/health/live",
}

var HealthReadyContract = contracts.RouteContract{
	Method: contracts.RouteMethodGet,
	Path:   "/health/ready",
}

type HealthRoutes struct {
	Live  contracts.RouteContract
	Ready contracts.RouteContract
}

var DefaultHealthRoutes = HealthRoutes{
	Live:  HealthLiveContract,
	Ready: HealthReadyContract,
}

var MetricsContract = contracts.RouteContract{
	Method: contracts.RouteMethodGet,
	Path:   "/metrics",
}

type MetricsRoutes struct {
	Collect contracts.RouteContract
}

var DefaultMetricsRoutes = MetricsRoutes{Collect: MetricsContract}
