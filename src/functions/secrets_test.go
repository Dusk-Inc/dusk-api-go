package functions

import (
	"reflect"
	"testing"
)

func TestDomain__ParseSecretsFile__ParsesValidEntries(t *testing.T) {
	content := "A=1\nexport B=2\n# comment\n"
	parsed := ParseSecretsFile(content)
	expected := map[string]string{"A": "1", "B": "2"}
	if !reflect.DeepEqual(parsed, expected) {
		t.Fatalf("expected %v, got %v", expected, parsed)
	}
}

func TestDomain__BuildRotation__ComputesDiffs(t *testing.T) {
	previous := map[string]string{"A": "1", "B": "2"}
	current := map[string]string{"B": "3", "C": "4"}
	rotation := BuildRotation(previous, current, 1, 2)
	if !reflect.DeepEqual(rotation.AddedKeys, []string{"C"}) {
		t.Fatalf("unexpected added keys: %v", rotation.AddedKeys)
	}
	if !reflect.DeepEqual(rotation.RemovedKeys, []string{"A"}) {
		t.Fatalf("unexpected removed keys: %v", rotation.RemovedKeys)
	}
	if !reflect.DeepEqual(rotation.UpdatedKeys, []string{"B"}) {
		t.Fatalf("unexpected updated keys: %v", rotation.UpdatedKeys)
	}
}
