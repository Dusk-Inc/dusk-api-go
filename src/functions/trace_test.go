package functions

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDomain__TraceMiddleware__PreservesProvidedCorrelationID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(TraceMiddleware)
	engine.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, GetCorrelationID(context))
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("x-correlation-id", "abc-123")
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	if response.Header().Get("x-correlation-id") != "abc-123" {
		t.Fatalf("expected correlation id header passthrough")
	}
}

func TestBoundary__TraceMiddleware__GeneratesCorrelationIDWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(TraceMiddleware)
	engine.GET("/", func(context *gin.Context) {
		context.String(http.StatusOK, GetCorrelationID(context))
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, request)

	value := response.Header().Get("x-correlation-id")
	if response.Code != http.StatusOK || value == "" {
		t.Fatalf("expected generated correlation id, got status %d and value %q", response.Code, value)
	}
}
