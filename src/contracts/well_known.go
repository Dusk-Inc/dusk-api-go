package contracts

type DiscoveryModel struct {
	ID   string   `json:"id"`
	Caps []string `json:"caps"`
}

type WellKnownRouterConfig struct {
	Issuer          string           `json:"issuer"`
	AvailableModels []DiscoveryModel `json:"availableModels,omitempty"`
	PublicKeySet    map[string]any   `json:"publicKeySet"`
}

var WellKnownOpenIDConfigurationContract = RouteContract{
	Method: RouteMethodGet,
	Path:   "/.well-known/openid-configuration",
}

var WellKnownJWKSContract = RouteContract{
	Method: RouteMethodGet,
	Path:   "/.well-known/jwks.json",
}

type WellKnownRoutes struct {
	OpenIDConfiguration RouteContract
	JWKS                RouteContract
}

var DefaultWellKnownRoutes = WellKnownRoutes{
	OpenIDConfiguration: WellKnownOpenIDConfigurationContract,
	JWKS:                WellKnownJWKSContract,
}
