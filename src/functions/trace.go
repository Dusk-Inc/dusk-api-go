package functions

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

const correlationIDHeader = "x-correlation-id"

func TraceMiddleware(context *gin.Context) {
	correlationID := context.GetHeader(correlationIDHeader)
	if correlationID == "" {
		correlationID = makeCorrelationID()
	}

	context.Header(correlationIDHeader, correlationID)
	SetCorrelationID(context, correlationID)
	context.Next()
}

func makeCorrelationID() string {
	bytes := make([]byte, 16)
	if _, readError := rand.Read(bytes); readError != nil {
		return "correlation-id-unavailable"
	}
	return hex.EncodeToString(bytes)
}
