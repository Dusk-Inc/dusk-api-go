package contracts

type ServiceDecoratorPhase string

const (
	ServiceDecoratorPhaseEncode ServiceDecoratorPhase = "encode"
	ServiceDecoratorPhaseDecode ServiceDecoratorPhase = "decode"
)

type ServiceDecoratorTransformErrorInput struct {
	Phase   ServiceDecoratorPhase
	Target  string
	Message string
}

type ServiceMapperContext struct {
	ServiceName string
	MethodName  string
	Phase       ServiceDecoratorPhase
}

type ServiceArgsMapper func(args []any, context ServiceMapperContext) ([]any, error)

type ServiceResultMapper func(result any, context ServiceMapperContext) (any, error)

type ServiceDecoratorRule struct {
	Methods   []string
	MapArgs   ServiceArgsMapper
	MapResult ServiceResultMapper
}

type ServiceDecoratorConfig struct {
	ServiceName string
	Rules       []ServiceDecoratorRule
}
