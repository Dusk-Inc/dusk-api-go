package modules

import (
	"fmt"
	"log/slog"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
)

type RuntimeManager struct {
	logger           *slog.Logger
	plugins          []contracts.RuntimePlugin
	startedPluginIDs []string
	dependencies     map[string]any
	started          bool
}

func NewRuntimeManager(logger *slog.Logger) *RuntimeManager {
	if logger == nil {
		logger = slog.Default()
	}
	return &RuntimeManager{
		logger:           logger,
		plugins:          make([]contracts.RuntimePlugin, 0),
		startedPluginIDs: make([]string, 0),
		dependencies:     map[string]any{},
	}
}

func (manager *RuntimeManager) Use(plugin contracts.RuntimePlugin) error {
	if manager.started {
		return fmt.Errorf("cannot register runtime plugin after startup")
	}
	for _, item := range manager.plugins {
		if item.ID() == plugin.ID() {
			return fmt.Errorf("runtime plugin already registered: %s", plugin.ID())
		}
	}
	manager.plugins = append(manager.plugins, plugin)
	return nil
}

func (manager *RuntimeManager) Start() error {
	if manager.started {
		return nil
	}
	context := manager.buildContext()
	for _, plugin := range manager.plugins {
		if setupError := plugin.Setup(context); setupError != nil {
			return setupError
		}
		if startError := plugin.Start(context); startError != nil {
			return startError
		}
		manager.startedPluginIDs = append(manager.startedPluginIDs, plugin.ID())
	}
	manager.started = true
	return nil
}

func (manager *RuntimeManager) Stop() error {
	if !manager.started {
		return nil
	}
	context := manager.buildContext()
	for index := len(manager.startedPluginIDs) - 1; index >= 0; index -= 1 {
		id := manager.startedPluginIDs[index]
		for _, plugin := range manager.plugins {
			if plugin.ID() != id {
				continue
			}
			if stopError := plugin.Stop(context); stopError != nil {
				return stopError
			}
		}
	}
	manager.startedPluginIDs = manager.startedPluginIDs[:0]
	manager.started = false
	return nil
}

func (manager *RuntimeManager) GetDependency(key string) any {
	return manager.dependencies[key]
}

func (manager *RuntimeManager) SetDependency(key string, value any) {
	manager.dependencies[key] = value
}

func (manager *RuntimeManager) buildContext() contracts.RuntimePluginContext {
	return contracts.RuntimePluginContext{
		Logger: manager.logger,
		SetDependency: func(key string, value any) {
			manager.SetDependency(key, value)
		},
		GetDependency: func(key string) any {
			return manager.GetDependency(key)
		},
	}
}
