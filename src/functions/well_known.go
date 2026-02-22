package functions

import (
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
)

func MakeOpenIDConfiguration(config contracts.WellKnownRouterConfig) map[string]any {
	payload := map[string]any{
		"issuer":                                config.Issuer,
		"jwks_uri":                              config.Issuer + contracts.DefaultWellKnownRoutes.JWKS.Path,
		"authorization_endpoint":                config.Issuer + "/authorize",
		"token_endpoint":                        config.Issuer + "/token",
		"userinfo_endpoint":                     config.Issuer + "/userinfo",
		"ai_endpoint":                           config.Issuer + "/ai/models",
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"response_types_supported":              []string{"code", "id_token"},
		"scopes_supported":                      []string{"openid", "profile", "email", "ai_access"},
	}

	if len(config.AvailableModels) > 0 {
		payload["ai_models_supported"] = config.AvailableModels
	}

	return payload
}
