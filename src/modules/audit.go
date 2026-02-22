package modules

import (
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/gin-gonic/gin"
)

func ReadTraceID(context *gin.Context) string {
	correlationID := context.GetHeader("x-correlation-id")
	if correlationID != "" {
		return correlationID
	}
	traceID := context.GetHeader("x-trace-id")
	if traceID != "" {
		return traceID
	}
	return ""
}

func AuditMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		traceID := ReadTraceID(context)
		context.Set("trace_id", traceID)
		context.Set("audit_log", func(level contracts.AuditLevel, payload contracts.AuditPayload) {
			if traceID != "" {
				payload["trace_id"] = traceID
			}
			_ = level
		})
		context.Next()
	}
}
