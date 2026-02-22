package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/gin-gonic/gin"
)

type fakeCollector struct{}

func (fakeCollector) Collect(_ context.Context) (string, []byte, error) {
	return "text/plain", []byte("ok_metric 1\n"), nil
}

func TestDomain__MetricsRouter__ServesCollectorPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	router := NewMetricsRouter(contracts.MetricsRouterConfig{Collector: fakeCollector{}})
	router.Register(engine)

	request := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	if response.Body.String() != "ok_metric 1\n" {
		t.Fatalf("unexpected body: %q", response.Body.String())
	}
}
