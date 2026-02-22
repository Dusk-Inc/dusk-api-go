package functions

import (
	"os"
	"testing"
)

func TestDomain__ParseEnv__ParsesHostAndPort(t *testing.T) {
	t.Setenv("HOST", "127.0.0.1")
	t.Setenv("PORT", "4000")
	env, err := ParseEnv()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if env.Host != "127.0.0.1" || env.Port != 4000 {
		t.Fatalf("unexpected env: %+v", env)
	}
}

func TestComplement__ParseEnv__RejectsInvalidPort(t *testing.T) {
	t.Setenv("HOST", "0.0.0.0")
	t.Setenv("PORT", "invalid")
	_, err := ParseEnv()
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestBoundary__ParseEnv__MissingHostUsesDefault(t *testing.T) {
	_ = os.Unsetenv("HOST")
	t.Setenv("PORT", "3000")
	env, err := ParseEnv()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if env.Host != "0.0.0.0" {
		t.Fatalf("expected default host, got %s", env.Host)
	}
}
