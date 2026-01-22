package mop

import (
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/solver"
)

type MopOption = solver.SolverOption

// WithEpsilonSteps sets the number of epsilon steps for the epsilon-constraint method.
func WithEpsilonSteps(steps int) MopOption {
	return func(cfg *common.SolverConfig) {
		cfg.EpsilonSteps = steps
	}
}

// WithWeightedSums sets the weight vectors for the weighted sum method.
func WithWeightedSums(weights [][]float64) MopOption {
	return func(cfg *common.SolverConfig) {
		cfg.WeightedSums = weights
	}
}

// with sets the entire SolverConfig (for internal use only).
func with(cfg_ common.SolverConfig) MopOption {
	return func(cfg *common.SolverConfig) {
		*cfg = cfg_
	}
}
