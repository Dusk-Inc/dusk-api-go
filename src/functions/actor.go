package functions

import (
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/gin-gonic/gin"
)

func ReadActorField(context *gin.Context, field string, source contracts.ActorSource) string {
	switch source {
	case contracts.ActorSourceHeader:
		return context.GetHeader(field)
	case contracts.ActorSourceQuery:
		return context.Query(field)
	case contracts.ActorSourceBody:
		var payload map[string]any
		if bindError := context.ShouldBindJSON(&payload); bindError != nil {
			return ""
		}
		value, ok := payload[field]
		if !ok {
			return ""
		}
		actorID, ok := value.(string)
		if !ok {
			return ""
		}
		return actorID
	default:
		return ""
	}
}

func MakeMissingActorPayload(code string, message string) contracts.ActorMiddlewareErrorResponse {
	return contracts.ActorMiddlewareErrorResponse{
		Error: contracts.ActorMiddlewareErrorBody{Code: code, Message: message},
	}
}

func SendMissingActor(context *gin.Context, payload contracts.ActorMiddlewareErrorResponse, statusCode int) {
	context.JSON(statusCode, payload)
}
