# Dusk API Go

Gin-based shared API runtime and utilities for Dusk Go services.

## Structure

- `src/modules`: runtime and middleware modules (`AppManager`, `RuntimeManager`, `SecretManager`, `ActorMiddleware`, `WellKnownRouter`, `ServiceDecorator`)
- `src/routes`: default routers (`HealthRouter`, `MetricsRouter`)
- `src/contracts`: shared contracts for runtime/plugins/routes/service decorators
- `src/functions`: pure helpers (env parsing, trace context, actor extraction, secrets parsing, well-known payload generation)
- `src/tokens`: constants for routes, actor defaults, service decorator, runtime dependency keys

## Quick Start

```go
package main

import (
    "log"

    duskapi "github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src"
    "github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
)

func main() {
    api := duskapi.NewAppManager(contracts.AppManagerConfig{ServiceName: "my-service"})

    env, err := duskapi.ParseEnv()
    if err != nil {
        log.Fatal(err)
    }

    if err := api.Engine.Run(env.Host + ":" + "3000"); err != nil {
        log.Fatal(err)
    }
}
```

## Runtime

- Runtime plugins are managed by `RuntimeManager`.
- Secrets can be attached with `UseSecrets(...)` and accessed through runtime dependencies.

## Validation

Run from `repos/libs/dusk-api-go`:

```bash
go test ./...
```
