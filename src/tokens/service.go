package tokens

import "github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"

const (
	ServiceDecoratorTransformErrorCode = "SERVICE_DECORATOR_TRANSFORM_ERROR"
	ServiceDecoratorSafeErrorMessage   = "Service transformation failed."
)

var ServiceDecoratorPhase = struct {
	Encode contracts.ServiceDecoratorPhase
	Decode contracts.ServiceDecoratorPhase
}{
	Encode: contracts.ServiceDecoratorPhaseEncode,
	Decode: contracts.ServiceDecoratorPhaseDecode,
}
