package mop

import (
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/solver"
)

type MopOption = solver.SolverOption

// WithEpsilonSteps sets the number of epsilon steps used by [SolvePareto]
// when the [ParetoEpsilonConstraint] method is selected.
func WithEpsilonSteps(steps int) MopOption {
	return func(cfg *common.SolverConfig) {
		cfg.EpsilonSteps = steps
	}
}

// WithWeightedSums sets weight vectors used by [SolvePareto] when the
// [ParetoWeightedSum] method is selected.
func WithWeightedSums(weights [][]float64) MopOption {
	return func(cfg *common.SolverConfig) {
		cfg.WeightedSums = weights
	}
}

// with sets the entire SolverConfig and is used internally to pass modified
// configs into internal solver calls.
func with(cfg_ common.SolverConfig) MopOption {
	return func(cfg *common.SolverConfig) {
		*cfg = cfg_
	}
}
