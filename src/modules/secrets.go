package modules

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/functions"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/tokens"
)

type SecretRotationListener func(rotation contracts.SecretRotation)

type SecretManager struct {
	env                 map[string]string
	secretPathEnvVar    string
	secretPathDefault   string
	requireReadOnlyFile bool
	listeners           []SecretRotationListener
	snapshot            contracts.SecretSnapshot
	lock                sync.RWMutex
}

func NewSecretManager(options contracts.SecretManagerOptions) *SecretManager {
	env := options.Env
	if env == nil {
		env = readProcessEnv()
	}
	secretPathEnvVar := options.SecretPathEnvVar
	if secretPathEnvVar == "" {
		secretPathEnvVar = tokens.DefaultSecretPathEnvVar
	}
	secretPathDefault := options.SecretPathDefault
	if secretPathDefault == "" {
		secretPathDefault = tokens.DefaultSecretPath
	}
	requireReadOnlyFile := true
	if options.RequireReadOnlyFile == false {
		requireReadOnlyFile = false
	}

	return &SecretManager{
		env:                 env,
		secretPathEnvVar:    secretPathEnvVar,
		secretPathDefault:   secretPathDefault,
		requireReadOnlyFile: requireReadOnlyFile,
		listeners:           make([]SecretRotationListener, 0),
		snapshot: contracts.SecretSnapshot{
			Generation: 0,
			Values:     map[string]string{},
		},
	}
}

func (manager *SecretManager) GetSnapshot() contracts.SecretSnapshot {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	copyValues := map[string]string{}
	for key, value := range manager.snapshot.Values {
		copyValues[key] = value
	}
	return contracts.SecretSnapshot{Generation: manager.snapshot.Generation, Values: copyValues}
}

func (manager *SecretManager) GetSecret(key string) string {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	return manager.snapshot.Values[key]
}

func (manager *SecretManager) GetRequiredSecret(key string) (string, error) {
	value := manager.GetSecret(key)
	if value == "" {
		return "", fmt.Errorf("required secret is missing: %s", key)
	}
	return value, nil
}

func (manager *SecretManager) GetAllSecrets() map[string]string {
	snapshot := manager.GetSnapshot()
	return snapshot.Values
}

func (manager *SecretManager) OnRotate(listener SecretRotationListener) func() {
	manager.lock.Lock()
	manager.listeners = append(manager.listeners, listener)
	index := len(manager.listeners) - 1
	manager.lock.Unlock()

	return func() {
		manager.lock.Lock()
		defer manager.lock.Unlock()
		if index >= 0 && index < len(manager.listeners) {
			manager.listeners[index] = nil
		}
	}
}

func (manager *SecretManager) LoadSecrets() (contracts.SecretSnapshot, error) {
	return manager.RefreshSecrets()
}

func (manager *SecretManager) EnsureFreshSecretsFile() error {
	path := functions.ResolveSecretPath(manager.env, manager.secretPathEnvVar, manager.secretPathDefault)
	stats, statError := os.Stat(path)
	if statError != nil {
		return statError
	}
	if time.Since(stats.ModTime()) > 5*time.Minute {
		return fmt.Errorf("secrets file is stale: %s", path)
	}
	return nil
}

func (manager *SecretManager) RefreshSecrets() (contracts.SecretSnapshot, error) {
	values, collectError := manager.collectSecrets()
	if collectError != nil {
		return contracts.SecretSnapshot{}, collectError
	}

	manager.lock.Lock()
	defer manager.lock.Unlock()
	if functions.AreSecretMapsEqual(manager.snapshot.Values, values) {
		return manager.snapshot, nil
	}

	previousValues := manager.snapshot.Values
	generation := manager.snapshot.Generation + 1
	manager.snapshot = contracts.SecretSnapshot{Generation: generation, Values: values}

	if generation > 1 {
		rotation := functions.BuildRotation(previousValues, values, generation-1, generation)
		for _, listener := range manager.listeners {
			if listener != nil {
				listener(rotation)
			}
		}
	}

	return manager.snapshot, nil
}

func (manager *SecretManager) StartWatching() error {
	return nil
}

func (manager *SecretManager) StopWatching() {}

func (manager *SecretManager) collectSecrets() (map[string]string, error) {
	secretPath := functions.ResolveSecretPath(manager.env, manager.secretPathEnvVar, manager.secretPathDefault)
	content, readError := os.ReadFile(secretPath)
	if readError != nil {
		if functions.IsMissingFileError(readError) {
			return manager.env, nil
		}
		return nil, readError
	}

	if manager.requireReadOnlyFile {
		if writable := fileIsWritable(secretPath); writable {
			return nil, fmt.Errorf("secrets file is writable by current process: %s", secretPath)
		}
	}

	fileSecrets := functions.ParseSecretsFile(string(content))
	return functions.MergeWithEnv(fileSecrets, manager.env), nil
}

func readProcessEnv() map[string]string {
	values := map[string]string{}
	for _, entry := range os.Environ() {
		separator := -1
		for index := range entry {
			if entry[index] == '=' {
				separator = index
				break
			}
		}
		if separator <= 0 {
			continue
		}
		values[entry[:separator]] = entry[separator+1:]
	}
	return values
}

func fileIsWritable(path string) bool {
	file, openError := os.OpenFile(path, os.O_WRONLY, 0)
	if openError != nil {
		return false
	}
	_ = file.Close()
	return true
}
