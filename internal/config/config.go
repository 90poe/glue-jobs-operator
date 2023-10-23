package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Configuration from Env
	OperatorConfig struct {
		// NOTE: Log Level is set via zap-log-level flag passed to the operator
		// MaxConcurrentReconciles is the maximum number of concurrent Reconciles which can be run.
		MaxConcurrentReconciles int `env:"MAX_CONCURRENT_RECONCILES" env-default:"1"`
	}
)

func New() (OperatorConfig, error) {
	// get config from Env
	cfg := OperatorConfig{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("can't make operator config: %w", err)
	}
	return cfg, nil
}
