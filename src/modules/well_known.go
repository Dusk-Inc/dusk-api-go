package modules

import (
	"net/http"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/functions"
	"github.com/gin-gonic/gin"
)

type WellKnownRouter struct {
	config contracts.WellKnownRouterConfig
}

func NewWellKnownRouter(config contracts.WellKnownRouterConfig) *WellKnownRouter {
	return &WellKnownRouter{config: config}
}

func (router *WellKnownRouter) Register(engine gin.IRoutes) {
	engine.GET(contracts.DefaultWellKnownRoutes.OpenIDConfiguration.Path, func(context *gin.Context) {
		context.JSON(http.StatusOK, functions.MakeOpenIDConfiguration(router.config))
	})
	engine.GET(contracts.DefaultWellKnownRoutes.JWKS.Path, func(context *gin.Context) {
		context.JSON(http.StatusOK, router.config.PublicKeySet)
	})
}
