package contracts

type AuditLevel string

const (
	AuditLevelInfo  AuditLevel = "info"
	AuditLevelError AuditLevel = "error"
)

type AuditPayload map[string]any

type RequestLogger interface {
	Info(payload AuditPayload)
	Error(payload AuditPayload)
}
