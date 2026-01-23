package solver

import (
	"context"

	"github.com/chriso345/gspl/internal/common"
)

// SolverOption is a functional option that modifies a SolverConfig.
type SolverOption func(*common.SolverConfig)

// WithTolerance sets the numerical tolerance used by the solver.
// It controls convergence criteria for linear solves and floating-point checks.
// It does not affect integer gap sensitivity; use [WithGapSensitivity] for that.
func WithTolerance(t float64) SolverOption {
	return func(cfg *common.SolverConfig) {
		cfg.Tolerance = t
	}
}

// WithContext attaches a context used to cancel or time out solver operations.
func WithContext(ctx context.Context) SolverOption {
	return func(cfg *common.SolverConfig) {
		cfg.Ctx = ctx
	}
}

// WithMaxIterations sets an upper bound on the number of solver iterations.
// A value <= 0 usually means no iteration limit.
func WithMaxIterations(max int) SolverOption {
	return func(cfg *common.SolverConfig) {
		cfg.MaxIterations = max
	}
}

// WithGapSensitivity sets the sensitivity used when comparing objective gaps
// in integer problems. Smaller values make the solver treat small differences
// as significant.
func WithGapSensitivity(gap float64) SolverOption {
	return func(cfg *common.SolverConfig) {
		cfg.GapSensitivity = gap
	}
}

// WithThreads configures the number of OS threads the solver may use.
// Currently not implemented and will panic if used.
func WithThreads(n int) SolverOption {
	panic("multi-threading not yet implemented")
}

// WithLogging enables or disables internal solver logging.
// Pass true to enable verbose runtime information useful for debugging.
func WithLogging(enabled bool) SolverOption {
	return func(cfg *common.SolverConfig) {
		cfg.Logging = enabled
	}
}

/// Strategy Functions Options

// WithBranch sets a custom branching strategy used by the solver.
// By default the solver uses [brancher.DefaultBranch] if no branch function is supplied.
func WithBranch(common.BranchFunc) SolverOption {
	panic("branching not yet implemented")
}

// WithHeuristic sets the heuristic function that produces incumbent solutions.
// If unspecified the solver uses [brancher.DefaultHeuristic].
func WithHeuristic(common.HeuristicFunc) SolverOption {
	panic("heuristics not yet implemented")
}

// WithCut sets the cut-generation function used to add cutting planes.
// The solver falls back to [brancher.DefaultCut] when none is provided.
func WithCut(common.CutFunc) SolverOption {
	panic("cutting planes not yet implemented")
}

/// Helpers

// NewSolverConfig returns a SolverConfig built from the default configuration
// with each provided SolverOption applied in order.
func NewSolverConfig(opts ...SolverOption) *common.SolverConfig {
	cfg := common.DefaultSolverConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
