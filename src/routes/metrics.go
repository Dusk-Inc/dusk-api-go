package routes

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/tokens"
	"github.com/gin-gonic/gin"
)

type defaultMetricsCollector struct{}

func (defaultMetricsCollector) Collect(_ context.Context) (string, []byte, error) {
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	payload := strings.Join([]string{
		"# TYPE go_mem_alloc_bytes gauge",
		fmt.Sprintf("go_mem_alloc_bytes %d", memory.Alloc),
		"# TYPE go_goroutines gauge",
		fmt.Sprintf("go_goroutines %d", runtime.NumGoroutine()),
	}, "\n") + "\n"
	return "text/plain; version=0.0.4", []byte(payload), nil
}

type MetricsRouter struct {
	collector contracts.MetricsCollector
}

func NewMetricsRouter(config contracts.MetricsRouterConfig) *MetricsRouter {
	collector := config.Collector
	if collector == nil {
		collector = defaultMetricsCollector{}
	}
	return &MetricsRouter{collector: collector}
}

func (router *MetricsRouter) Register(engine gin.IRoutes) {
	engine.GET(tokens.DefaultMetricsRoutes.Collect.Path, func(context *gin.Context) {
		contentType, payload, collectError := router.collector.Collect(context.Request.Context())
		if collectError != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": "failed to collect metrics"}})
			return
		}
		context.Header("Content-Type", contentType)
		context.Data(http.StatusOK, contentType, payload)
	})
}
