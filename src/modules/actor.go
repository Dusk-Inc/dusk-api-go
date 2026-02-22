package modules

import (
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/functions"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/tokens"
	"github.com/gin-gonic/gin"
)

type ActorMiddleware struct {
	field             string
	source            contracts.ActorSource
	required          bool
	missingStatusCode int
	missingCode       string
	missingMessage    string
	readActor         contracts.ActorReader
	onMissingActor    contracts.MissingActorHandler
}

func NewActorMiddleware(field string) *ActorMiddleware {
	return &ActorMiddleware{
		field:             field,
		source:            tokens.ActorDefaultSource,
		required:          tokens.ActorDefaultRequired,
		missingStatusCode: tokens.ActorDefaultMissingStatusCode,
		missingCode:       tokens.ActorDefaultMissingCode,
		missingMessage:    tokens.ActorDefaultMissingMessage,
		readActor:         functions.ReadActorField,
		onMissingActor:    functions.SendMissingActor,
	}
}

func (middleware *ActorMiddleware) Handler() gin.HandlerFunc {
	return func(context *gin.Context) {
		actorID := middleware.readActor(context, middleware.field, middleware.source)
		if actorID == "" {
			if !middleware.required {
				context.Next()
				return
			}
			payload := functions.MakeMissingActorPayload(middleware.missingCode, middleware.missingMessage)
			middleware.onMissingActor(context, payload, middleware.missingStatusCode)
			context.Abort()
			return
		}
		context.Set("actor_id", actorID)
		context.Next()
	}
}
