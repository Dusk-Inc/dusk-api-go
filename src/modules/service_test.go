package modules

import (
	"errors"
	"testing"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
)

func TestDomain__ServiceDecorator__MapsArgsAndResult(t *testing.T) {
	decorator := NewServiceDecorator(contracts.ServiceDecoratorConfig{
		ServiceName: "service",
		Rules: []contracts.ServiceDecoratorRule{
			{
				MapArgs: func(args []any, _ contracts.ServiceMapperContext) ([]any, error) {
					return []any{args[0].(int) + 1}, nil
				},
				MapResult: func(result any, _ contracts.ServiceMapperContext) (any, error) {
					return result.(int) + 1, nil
				},
			},
		},
	})

	result, err := decorator.MapCall("save", []any{1}, func(args []any) (any, error) {
		return args[0].(int), nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.(int) != 3 {
		t.Fatalf("expected 3, got %v", result)
	}
}

func TestChaos__ServiceDecorator__WrapsMapErrors(t *testing.T) {
	decorator := NewServiceDecorator(contracts.ServiceDecoratorConfig{
		ServiceName: "service",
		Rules: []contracts.ServiceDecoratorRule{
			{MapArgs: func(_ []any, _ contracts.ServiceMapperContext) ([]any, error) {
				return nil, errors.New("boom")
			}},
		},
	})

	_, err := decorator.MapCall("save", []any{1}, func(args []any) (any, error) { return args, nil })
	if err == nil {
		t.Fatal("expected transform error")
	}
	if _, ok := err.(*ServiceDecoratorTransformError); !ok {
		t.Fatalf("expected ServiceDecoratorTransformError, got %T", err)
	}
}
