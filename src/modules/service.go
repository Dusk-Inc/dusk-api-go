package modules

import (
	"fmt"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/tokens"
)

type ServiceDecoratorTransformError struct {
	Code    string
	Phase   contracts.ServiceDecoratorPhase
	Target  string
	Message string
}

func (err *ServiceDecoratorTransformError) Error() string {
	return err.Message
}

type ServiceDecorator struct {
	serviceName string
	rules       []contracts.ServiceDecoratorRule
}

func NewServiceDecorator(config contracts.ServiceDecoratorConfig) *ServiceDecorator {
	serviceName := config.ServiceName
	if serviceName == "" {
		serviceName = "service"
	}
	return &ServiceDecorator{serviceName: serviceName, rules: config.Rules}
}

func (decorator *ServiceDecorator) MapCall(methodName string, args []any, call func([]any) (any, error)) (any, error) {
	mappedArgs, mapArgsError := decorator.mapCallArgs(args, methodName)
	if mapArgsError != nil {
		return nil, mapArgsError
	}

	result, callError := call(mappedArgs)
	if callError != nil {
		return nil, callError
	}

	return decorator.mapCallResult(result, methodName)
}

func (decorator *ServiceDecorator) mapCallArgs(args []any, methodName string) ([]any, error) {
	nextArgs := args
	for _, rule := range decorator.rules {
		if !shouldApplyRule(rule, methodName) || rule.MapArgs == nil {
			continue
		}
		context := contracts.ServiceMapperContext{
			ServiceName: decorator.serviceName,
			MethodName:  methodName,
			Phase:       contracts.ServiceDecoratorPhaseEncode,
		}
		mapped, mapError := rule.MapArgs(nextArgs, context)
		if mapError != nil {
			return nil, wrapTransformError(mapError, context)
		}
		nextArgs = mapped
	}
	return nextArgs, nil
}

func (decorator *ServiceDecorator) mapCallResult(result any, methodName string) (any, error) {
	nextResult := result
	for _, rule := range decorator.rules {
		if !shouldApplyRule(rule, methodName) || rule.MapResult == nil {
			continue
		}
		context := contracts.ServiceMapperContext{
			ServiceName: decorator.serviceName,
			MethodName:  methodName,
			Phase:       contracts.ServiceDecoratorPhaseDecode,
		}
		mapped, mapError := rule.MapResult(nextResult, context)
		if mapError != nil {
			return nil, wrapTransformError(mapError, context)
		}
		nextResult = mapped
	}
	return nextResult, nil
}

func shouldApplyRule(rule contracts.ServiceDecoratorRule, methodName string) bool {
	if len(rule.Methods) == 0 {
		return true
	}
	for _, method := range rule.Methods {
		if method == methodName {
			return true
		}
	}
	return false
}

func wrapTransformError(_ error, context contracts.ServiceMapperContext) *ServiceDecoratorTransformError {
	return &ServiceDecoratorTransformError{
		Code:    tokens.ServiceDecoratorTransformErrorCode,
		Phase:   context.Phase,
		Target:  fmt.Sprintf("%s.%s", context.ServiceName, context.MethodName),
		Message: tokens.ServiceDecoratorSafeErrorMessage,
	}
}
