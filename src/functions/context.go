package functions

import "github.com/gin-gonic/gin"

const correlationIDContextKey = "dusk.correlation_id"

func SetCorrelationID(context *gin.Context, correlationID string) {
	context.Set(correlationIDContextKey, correlationID)
}

func GetCorrelationID(context *gin.Context) string {
	value, exists := context.Get(correlationIDContextKey)
	if !exists {
		return "no-context"
	}
	correlationID, ok := value.(string)
	if !ok || correlationID == "" {
		return "no-context"
	}
	return correlationID
}
