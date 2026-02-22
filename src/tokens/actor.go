package tokens

import "github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"

const (
	ActorDefaultRequired          = true
	ActorDefaultMissingStatusCode = 401
	ActorDefaultMissingCode       = "MISSING_ACTOR"
	ActorDefaultMissingMessage    = "Actor is required."
)

var ActorDefaultSource = contracts.ActorSourceHeader
