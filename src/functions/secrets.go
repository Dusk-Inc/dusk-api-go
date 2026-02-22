package functions

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
)

func ParseSecretLine(line string) (string, string, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return "", "", false
	}

	normalized := strings.TrimPrefix(trimmed, "export ")
	separator := strings.Index(normalized, "=")
	if separator <= 0 {
		return "", "", false
	}

	key := strings.TrimSpace(normalized[:separator])
	if key == "" {
		return "", "", false
	}

	value := strings.TrimSpace(normalized[separator+1:])
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) >= 2 {
		value = strings.ReplaceAll(value[1:len(value)-1], "\\n", "\n")
	}

	return key, value, true
}

func ParseSecretsFile(content string) map[string]string {
	parsed := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		key, value, ok := ParseSecretLine(scanner.Text())
		if !ok {
			continue
		}
		parsed[key] = value
	}
	return parsed
}

func MergeWithEnv(fileSecrets map[string]string, env map[string]string) map[string]string {
	merged := map[string]string{}
	for key, value := range fileSecrets {
		merged[key] = value
	}
	for key, value := range env {
		merged[key] = value
	}
	return merged
}

func AreSecretMapsEqual(left map[string]string, right map[string]string) bool {
	if len(left) != len(right) {
		return false
	}
	for key, value := range left {
		if right[key] != value {
			return false
		}
	}
	return true
}

func BuildRotation(previousValues map[string]string, currentValues map[string]string, previousGeneration int, generation int) contracts.SecretRotation {
	previousKeys := map[string]struct{}{}
	for key := range previousValues {
		previousKeys[key] = struct{}{}
	}
	currentKeys := map[string]struct{}{}
	for key := range currentValues {
		currentKeys[key] = struct{}{}
	}

	addedKeys := make([]string, 0)
	removedKeys := make([]string, 0)
	updatedKeys := make([]string, 0)
	unchangedKeys := make([]string, 0)

	for key := range currentKeys {
		if _, ok := previousKeys[key]; !ok {
			addedKeys = append(addedKeys, key)
			continue
		}
		if previousValues[key] == currentValues[key] {
			unchangedKeys = append(unchangedKeys, key)
		} else {
			updatedKeys = append(updatedKeys, key)
		}
	}

	for key := range previousKeys {
		if _, ok := currentKeys[key]; !ok {
			removedKeys = append(removedKeys, key)
		}
	}

	sort.Strings(addedKeys)
	sort.Strings(removedKeys)
	sort.Strings(updatedKeys)
	sort.Strings(unchangedKeys)

	return contracts.SecretRotation{
		Generation:         generation,
		PreviousGeneration: previousGeneration,
		AddedKeys:          addedKeys,
		RemovedKeys:        removedKeys,
		UpdatedKeys:        updatedKeys,
		UnchangedKeys:      unchangedKeys,
	}
}

func ResolveSecretPath(env map[string]string, secretPathEnvVar string, secretPathDefault string) string {
	if env != nil {
		if configuredPath, ok := env[secretPathEnvVar]; ok && configuredPath != "" {
			return filepath.Clean(configuredPath)
		}
	}
	return filepath.Clean(secretPathDefault)
}

func IsMissingFileError(err error) bool {
	return errors.Is(err, os.ErrNotExist)
}
