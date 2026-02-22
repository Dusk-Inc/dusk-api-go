package tokens

import "time"

const (
	DefaultSecretPathEnvVar      = "DUSK_SECRETS_FILE"
	DefaultSecretPath            = "/var/run/secrets/dusk/secrets.env"
	RuntimePluginSecrets         = "secrets"
	RuntimeDependencySecretsMgr  = "runtime.secrets.manager"
	RuntimeDependencySecretsSnap = "runtime.secrets.snapshot"
	RuntimeDependencySecretsEnv  = "runtime.secrets.env"
)

var DefaultWatchDebounce = 250 * time.Millisecond
