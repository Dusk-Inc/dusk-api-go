package modules

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/gin-gonic/gin"
)

func TestDomain__WellKnownRouter__ServesOpenIDAndJWKS(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	router := NewWellKnownRouter(contracts.WellKnownRouterConfig{Issuer: "https://issuer.example.com", PublicKeySet: map[string]any{"keys": []any{}}})
	router.Register(engine)

	for _, path := range []string{"/.well-known/openid-configuration", "/.well-known/jwks.json"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("expected 200 for %s, got %d", path, response.Code)
		}
	}
}
