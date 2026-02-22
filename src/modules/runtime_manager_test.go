package modules

import (
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"testing"
)

type lifecyclePlugin struct {
	id      string
	started bool
	stopped bool
}

func (plugin *lifecyclePlugin) ID() string                                   { return plugin.id }
func (plugin *lifecyclePlugin) Setup(_ contracts.RuntimePluginContext) error { return nil }
func (plugin *lifecyclePlugin) Start(_ contracts.RuntimePluginContext) error {
	plugin.started = true
	return nil
}
func (plugin *lifecyclePlugin) Stop(_ contracts.RuntimePluginContext) error {
	plugin.stopped = true
	return nil
}

func TestDomain__RuntimeManager__StartsAndStopsPlugins(t *testing.T) {
	manager := NewRuntimeManager(nil)
	plugin := &lifecyclePlugin{id: "a"}
	if err := manager.Use(plugin); err != nil {
		t.Fatalf("unexpected use error: %v", err)
	}
	if err := manager.Start(); err != nil {
		t.Fatalf("unexpected start error: %v", err)
	}
	if err := manager.Stop(); err != nil {
		t.Fatalf("unexpected stop error: %v", err)
	}
	if !plugin.started || !plugin.stopped {
		t.Fatalf("expected lifecycle events, got started=%v stopped=%v", plugin.started, plugin.stopped)
	}
}

func TestComplement__RuntimeManager__RejectsDuplicatePluginIDs(t *testing.T) {
	manager := NewRuntimeManager(nil)
	if err := manager.Use(&lifecyclePlugin{id: "a"}); err != nil {
		t.Fatalf("unexpected first use error: %v", err)
	}
	if err := manager.Use(&lifecyclePlugin{id: "a"}); err == nil {
		t.Fatal("expected duplicate plugin error")
	}
}
