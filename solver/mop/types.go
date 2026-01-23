package mop

import (
	"github.com/chriso345/gspl/internal/common"
	"gonum.org/v1/gonum/mat"
)

// MopSolution contains the result of a multi-objective optimization.
//
// The [MopSolution] holds objective values (in the same order as provided
// in the program), the primal solution vector, and the [common.SolverStatus].
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
