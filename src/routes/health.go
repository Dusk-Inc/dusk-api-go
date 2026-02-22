package routes

import (
	"context"
	"net/http"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/tokens"
	"github.com/gin-gonic/gin"
)

type HealthRouter struct {
	readiness contracts.ReadinessCheck
}

func NewHealthRouter(config contracts.HealthRouterConfig) *HealthRouter {
	readiness := config.Readiness
	if readiness == nil {
		readiness = func(_ context.Context) bool { return true }
	}
	return &HealthRouter{readiness: readiness}
}

func (router *HealthRouter) Register(engine gin.IRoutes) {
	okPayload := gin.H{"data": gin.H{"status": "ok"}}
	engine.GET(tokens.DefaultHealthRoutes.Live.Path, func(context *gin.Context) {
		context.JSON(http.StatusOK, okPayload)
	})
	engine.GET(tokens.DefaultHealthRoutes.Ready.Path, func(context *gin.Context) {
		if router.readiness(context.Request.Context()) {
			context.JSON(http.StatusOK, okPayload)
			return
		}
		context.JSON(http.StatusServiceUnavailable, gin.H{"data": gin.H{"status": "unready"}})
	})
}
