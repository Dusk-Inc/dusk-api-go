package contracts

import "github.com/gin-gonic/gin"

type ActorSource string

const (
	ActorSourceHeader ActorSource = "header"
	ActorSourceQuery  ActorSource = "query"
	ActorSourceBody   ActorSource = "body"
)

type RequestData map[string]any

type ActorMiddlewareErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ActorMiddlewareErrorResponse struct {
	Error ActorMiddlewareErrorBody `json:"error"`
}

type ActorReader func(context *gin.Context, field string, source ActorSource) string

type MissingActorHandler func(context *gin.Context, payload ActorMiddlewareErrorResponse, statusCode int)
