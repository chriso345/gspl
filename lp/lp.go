package lp

import (
	"github.com/chriso345/gspl/internal/common"
	"gonum.org/v1/gonum/mat"
)

// LinearProgram represents a linear programming problem in standard form.
type LinearProgram struct {
	// Problem definition
	Objective   *mat.VecDense      // c
	Constraints *mat.Dense         // A
	RHS         *mat.VecDense      // b
	Sense       LpSense            // Minimize or Maximize
	ConTypes    []LpConstraintType // metadata for constraints
	Vars        []LpVariable       // metadata for variables

	// Solution
	ObjectiveValue float64
	PrimalSolution *mat.VecDense // x*
	// DualSolution   *mat.VecDense // y*
	Status common.SolverStatus

	// Multi-Objective support
	SecondaryObjectives []*mat.VecDense // c2, c3, ...

	// Metadata
	Description string

	// Internal flag indicating the stored Objective has been negated to represent
	// a maximization problem.
	ObjectiveIsNegated bool
}

// NewLinearProgram Create a new Linear Program
func NewLinearProgram(desc string, vars []LpVariable) LinearProgram {

	lp := LinearProgram{
		Description: desc,
		Vars:        make([]LpVariable, len(vars)),
	}

	copy(lp.Vars, vars)

	lp.Status = common.SolverStatusNotSolved
	return lp
}
