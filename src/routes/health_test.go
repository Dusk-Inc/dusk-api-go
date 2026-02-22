package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/gin-gonic/gin"
)

func TestDomain__HealthRouter__LiveAndReadyOK(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	router := NewHealthRouter(contracts.HealthRouterConfig{Readiness: func(_ context.Context) bool { return true }})
	router.Register(engine)

	for _, path := range []string{"/health/live", "/health/ready"} {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()
		engine.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("expected 200 for %s, got %d", path, response.Code)
		}
	}
}

func TestComplement__HealthRouter__ReadyUnready503(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	router := NewHealthRouter(contracts.HealthRouterConfig{Readiness: func(_ context.Context) bool { return false }})
	router.Register(engine)

	request := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", response.Code)
	}
}
