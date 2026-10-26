package common

import (
	"github.com/chriso345/gspl/internal/matrix"
	"gonum.org/v1/gonum/mat"
)

// StandardComputationalForm represents a linear programming problem in standard form.
type StandardComputationalForm struct {
	Objective   *mat.VecDense // c
	Constraints *mat.Dense    // A
	RHS         *mat.VecDense // b

	PrimalSolution *mat.VecDense // x*

	ObjectiveValue *float64
	Status         *SolverStatus // Optimal, Infeasible, Unbounded, etc.
	SlackIndices   []int         // Indices of slack variables in the solution
	NumPrimals     int           // Number of primal variables (non-slack)

	IsMaximization bool

	VarCategories []VarCategory
}

// Copy creates a deep copy of the SCF
func (scf *StandardComputationalForm) Copy() *StandardComputationalForm {
	// Deep-copy pointer fields to avoid sharing mutable state between SCFs
	var objValPtr *float64
	if scf.ObjectiveValue != nil {
		v := *scf.ObjectiveValue
		objValPtr = new(float64)
		*objValPtr = v
	}
	var statusPtr *SolverStatus
	if scf.Status != nil {
		s := *scf.Status
		statusPtr = new(SolverStatus)
		*statusPtr = s
	}
	// Copy slack indices slice
	slackCopy := make([]int, len(scf.SlackIndices))
	copy(slackCopy, scf.SlackIndices)

	return &StandardComputationalForm{
		Objective:      mat.VecDenseCopyOf(scf.Objective),
		Constraints:    mat.DenseCopyOf(scf.Constraints),
		RHS:            mat.VecDenseCopyOf(scf.RHS),
		PrimalSolution: mat.VecDenseCopyOf(scf.PrimalSolution),
		ObjectiveValue: objValPtr,
		Status:         statusPtr,
		SlackIndices:   slackCopy,
		NumPrimals:     scf.NumPrimals,
		IsMaximization: scf.IsMaximization,
		VarCategories:  append([]VarCategory(nil), scf.VarCategories...),
	}
}

// AddBranch adds the bound x_idx <= rhs (dir 1) or x_idx >= rhs (dir 2) to the SCF.
//
// Rows in the SCF are equalities, so the bound gets its own slack (dir 1) or
// surplus (dir 2) column. Without it the branch would fix x_idx to rhs.
func (scf *StandardComputationalForm) AddBranch(idx int, rhs float64, dir int) {
	numRows, numCols := scf.Constraints.Dims()
	newConstraints := matrix.ResizeMatDense(scf.Constraints, numRows+1, numCols+1)
	newConstraints.Set(numRows, idx, 1)
	switch dir {
	case 1:
		newConstraints.Set(numRows, numCols, 1) // slack
	case 2:
		newConstraints.Set(numRows, numCols, -1) // surplus
	}
	newRHS := matrix.ResizeVecDense(scf.RHS, numRows+1)
	newRHS.SetVec(numRows, rhs)

	scf.Constraints = newConstraints
	scf.RHS = newRHS
	// The new column has no cost and is tracked like any other slack
	scf.Objective = matrix.ResizeVecDense(scf.Objective, numCols+1)
	scf.SlackIndices = append(scf.SlackIndices, numCols)
	if scf.VarCategories != nil {
		scf.VarCategories = append(scf.VarCategories, VarCategoryContinuous)
	}
}

// AddEquality appends an equality constraint fixing column idx to value.
func (scf *StandardComputationalForm) AddEquality(idx int, value float64) {
	numRows, numCols := scf.Constraints.Dims()
	newConstraints := mat.NewDense(numRows+1, numCols, nil)
	for r := range numRows {
		for c := range numCols {
			newConstraints.Set(r, c, scf.Constraints.At(r, c))
		}
	}
	for c := range numCols {
		if c == idx {
			newConstraints.Set(numRows, c, 1)
		} else {
			newConstraints.Set(numRows, c, 0)
		}
	}
	newRHS := mat.NewVecDense(numRows+1, nil)
	for r := range numRows {
		newRHS.SetVec(r, scf.RHS.AtVec(r))
	}
	newRHS.SetVec(numRows, value)
	scf.Constraints = newConstraints
	scf.RHS = newRHS
}
