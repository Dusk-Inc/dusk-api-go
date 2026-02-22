package contracts

import "log/slog"

type RuntimePluginContext struct {
	Logger        *slog.Logger
	SetDependency func(key string, value any)
	GetDependency func(key string) any
}

type RuntimePlugin interface {
	ID() string
	Setup(context RuntimePluginContext) error
	Start(context RuntimePluginContext) error
	Stop(context RuntimePluginContext) error
}
