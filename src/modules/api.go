package modules

import (
	"log/slog"
	"os"

	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/contracts"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/functions"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/routes"
	"github.com/dusk-inc/dusk-ocean/repos/libs/dusk-api-go/src/tokens"
	"github.com/gin-gonic/gin"
)

type AppManager struct {
	Engine  *gin.Engine
	Logger  *slog.Logger
	Runtime *RuntimeManager
}

func NewAppManager(config contracts.AppManagerConfig) *AppManager {
	logger := config.Logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	}

	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(functions.TraceMiddleware)
	engine.Use(AuditMiddleware())

	healthRouter := routes.NewHealthRouter(contracts.HealthRouterConfig{Readiness: config.Readiness})
	healthRouter.Register(engine)
	metricsRouter := routes.NewMetricsRouter(contracts.MetricsRouterConfig{})
	metricsRouter.Register(engine)

	return &AppManager{Engine: engine, Logger: logger, Runtime: NewRuntimeManager(logger)}
}

func (manager *AppManager) Use(plugin contracts.RuntimePlugin) error {
	return manager.Runtime.Use(plugin)
}

func (manager *AppManager) UseSecrets(config contracts.SecretManagerOptions) error {
	return manager.Use(NewSecretsPlugin(config))
}

func (manager *AppManager) GetDependency(key string) any {
	return manager.Runtime.GetDependency(key)
}

func (manager *AppManager) GetSecretsManager() *SecretManager {
	value := manager.GetDependency(tokens.RuntimeDependencySecretsMgr)
	secretManager, _ := value.(*SecretManager)
	return secretManager
}

func (manager *AppManager) GetSecretsSnapshot() contracts.SecretSnapshot {
	value := manager.GetDependency(tokens.RuntimeDependencySecretsSnap)
	snapshot, ok := value.(contracts.SecretSnapshot)
	if !ok {
		return contracts.SecretSnapshot{Generation: 0, Values: map[string]string{}}
	}
	return snapshot
}

func (manager *AppManager) GetSecretsEnv() map[string]string {
	value := manager.GetDependency(tokens.RuntimeDependencySecretsEnv)
	env, ok := value.(map[string]string)
	if !ok {
		return map[string]string{}
	}
	return env
}

func (manager *AppManager) StartRuntime() error {
	return manager.Runtime.Start()
}

func (manager *AppManager) StopRuntime() error {
	return manager.Runtime.Stop()
}
