package contracts

import "time"

type SecretSnapshot struct {
	Generation int
	Values     map[string]string
}

type SecretRotation struct {
	Generation         int
	PreviousGeneration int
	AddedKeys          []string
	RemovedKeys        []string
	UpdatedKeys        []string
	UnchangedKeys      []string
}

type SecretManagerOptions struct {
	Env                 map[string]string
	SecretPathEnvVar    string
	SecretPathDefault   string
	WatchDebounce       time.Duration
	RequireReadOnlyFile bool
}

type SecretLogger interface {
	Info(message string, attrs map[string]any)
	Warn(message string, attrs map[string]any)
}
