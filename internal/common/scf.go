package common

import (
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

// AddBranch adds a new constraint to the SCF
func (scf *StandardComputationalForm) AddBranch(idx int, rhs float64, dir int) {
	numRows, numCols := scf.Constraints.Dims()
	newConstraints := mat.NewDense(numRows+1, numCols, nil)
	newRHS := mat.NewVecDense(numRows+1, nil)
	for i := range numRows {
		for j := range numCols {
			newConstraints.Set(i, j, scf.Constraints.At(i, j))
		}
		newRHS.SetVec(i, scf.RHS.AtVec(i))
	}

	for j := range numCols {
		if j == idx {
			switch dir {
			case 1:
				newConstraints.Set(numRows, j, 1)
			case 2:
				newConstraints.Set(numRows, j, -1)
			}
		} else {
			newConstraints.Set(numRows, j, 0)
		}
	}
	switch dir {
	case 1:
		newRHS.SetVec(numRows, rhs)
	case 2:
		newRHS.SetVec(numRows, -rhs)
	}
	scf.Constraints = newConstraints
	scf.RHS = newRHS
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
