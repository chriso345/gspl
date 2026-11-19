package common

import (
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestDefaultSolverConfig(t *testing.T) {
	cfg := DefaultSolverConfig()
	assert.False(t, cfg.Logging)
	assert.Equal(t, cfg.Tolerance, 1e-6)
	assert.Equal(t, cfg.MaxIterations, 1000)
	assert.Equal(t, cfg.GapSensitivity, 0.05)
	assert.True(t, cfg.Branch == nil)
	assert.True(t, cfg.Heuristic == nil)
	assert.True(t, cfg.Cut == nil)
	assert.Equal(t, cfg.Threads, 0)
}

func TestValidateSolverConfig(t *testing.T) {
	cfg := DefaultSolverConfig()
	err := ValidateSolverConfig(cfg)
	assert.Nil(t, err)
}

func TestValidateSolverConfig_NonPositiveIterationLimit(t *testing.T) {
	cfg := DefaultSolverConfig()

	cfg.MaxIterations = -1
	err := ValidateSolverConfig(cfg)
	assert.NotNil(t, err)

	cfg.MaxIterations = 0
	err = ValidateSolverConfig(cfg)
	assert.Nil(t, err)
}
