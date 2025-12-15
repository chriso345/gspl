package common

import (
	"context"
	"github.com/chriso345/gspl/internal/errors"
)

// SolverConfig holds the actual configuration with no pointers.
type SolverConfig struct {
	Logging       bool
	Tolerance     float64
	MaxIterations int

	// Context for cancellation
	Ctx context.Context

	// IP Specific Options
	GapSensitivity float64
	Branch         BranchFunc
	Heuristic      HeuristicFunc
	Cut            CutFunc

	// Not Yet Implemented
	Threads int

	Debug bool
}

// DefaultSolverConfig returns the default solver configuration.
func DefaultSolverConfig() *SolverConfig {
	return &SolverConfig{
		Logging:       false, // Default logging is off
		Tolerance:     1e-6,
		MaxIterations: 1000,
		Ctx:           context.Background(),

		GapSensitivity: 0.05,
		Branch:         nil, // Default branching strategy defined in `brancher`
		Heuristic:      nil, // Default heuristic defined in `brancher`
		Cut:            nil, // Default cutting planes defined in `brancher`

		Threads: 0, // 0 means use all available cores

		Debug: false,
	}
}

// ValidateSolverConfig checks if the SolverConfig is valid.
func ValidateSolverConfig(cfg *SolverConfig) error {
	if cfg == nil {
		return errors.New(errors.ErrInvalidInput, "solver config is nil", nil)
	}
	if cfg.Tolerance <= 0 {
		return errors.New(errors.ErrInvalidInput, "tolerance must be > 0", nil)
	}
	if cfg.MaxIterations <= 0 {
		return errors.New(errors.ErrInvalidInput, "max iterations must be > 0", nil)
	}
	if cfg.GapSensitivity < 0 || cfg.GapSensitivity > 1 {
		return errors.New(errors.ErrInvalidInput, "gap sensitivity must be between 0 and 1", nil)
	}

	if cfg.Debug {
		cfg.Logging = true
	}

	if cfg.Ctx == nil {
		cfg.Ctx = context.Background()
	}

	return nil
}
