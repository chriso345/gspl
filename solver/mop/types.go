package mop

import (
	"github.com/chriso345/gspl/internal/common"
	"gonum.org/v1/gonum/mat"
)

// MopSolution contains the result of a multi-objective optimization.
type MopSolution struct {
	ObjectiveValues []float64
	PrimalSolution  *mat.VecDense
	Status          common.SolverStatus
}

type ParetoMethod int

const (
	ParetoWeightedSum ParetoMethod = iota
	ParetoEpsilonConstraint
)
