package modules

import (
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/tokens"
)

type SecretsPlugin struct {
	config      contracts.SecretManagerOptions
	manager     *SecretManager
	unsubscribe func()
}

func NewSecretsPlugin(config contracts.SecretManagerOptions) *SecretsPlugin {
	return &SecretsPlugin{config: config}
}

func (plugin *SecretsPlugin) ID() string {
	return tokens.RuntimePluginSecrets
}

func (plugin *SecretsPlugin) Setup(_ contracts.RuntimePluginContext) error {
	return nil
}

func (plugin *SecretsPlugin) Start(context contracts.RuntimePluginContext) error {
	plugin.manager = NewSecretManager(plugin.config)
	snapshot, loadError := plugin.manager.LoadSecrets()
	if loadError != nil {
		return loadError
	}

	context.SetDependency(tokens.RuntimeDependencySecretsMgr, plugin.manager)
	context.SetDependency(tokens.RuntimeDependencySecretsSnap, snapshot)
	context.SetDependency(tokens.RuntimeDependencySecretsEnv, snapshot.Values)

	plugin.unsubscribe = plugin.manager.OnRotate(func(rotation contracts.SecretRotation) {
		if plugin.manager == nil {
			return
		}
		latestSnapshot := plugin.manager.GetSnapshot()
		context.SetDependency(tokens.RuntimeDependencySecretsSnap, latestSnapshot)
		context.SetDependency(tokens.RuntimeDependencySecretsEnv, latestSnapshot.Values)
		if context.Logger != nil {
			context.Logger.Info("secret rotation detected", "generation", rotation.Generation, "previous_generation", rotation.PreviousGeneration)
		}
	})

	return plugin.manager.StartWatching()
}

func (plugin *SecretsPlugin) Stop(_ contracts.RuntimePluginContext) error {
	if plugin.unsubscribe != nil {
		plugin.unsubscribe()
		plugin.unsubscribe = nil
	}
	if plugin.manager != nil {
		plugin.manager.StopWatching()
		plugin.manager = nil
	}
	return nil
}
