package common

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"gonum.org/v1/gonum/mat"
)

func TestSCFCopy(t *testing.T) {
	obj := mat.NewVecDense(2, []float64{1, 2})
	constr := mat.NewDense(1, 2, []float64{3, 4})
	rhs := mat.NewVecDense(1, []float64{5})
	primal := mat.NewVecDense(2, []float64{6, 7})
	objVal := 10.0
	status := SolverStatusOptimal

	scf := &StandardComputationalForm{
		Objective:      obj,
		Constraints:    constr,
		RHS:            rhs,
		PrimalSolution: primal,
		ObjectiveValue: &objVal,
		Status:         &status,
		SlackIndices:   []int{0, 1},
		NumPrimals:     2,
	}

	copySCF := scf.Copy()
	// Ensure dimensions and values match
	assert.Equal(t, copySCF.Objective.Len(), 2)
	assert.Equal(t, copySCF.Constraints.RawMatrix().Rows, 1)
	assert.Equal(t, copySCF.RHS.Len(), 1)
	assert.Equal(t, *copySCF.ObjectiveValue, objVal)
	assert.Equal(t, *copySCF.Status, status)
	// Ensure deep copy: modifying original doesn't affect copy
	obj.SetVec(0, 999)
	assert.Equal(t, copySCF.Objective.AtVec(0), 1.0)
}

func TestSCFAddBranch(t *testing.T) {
	obj := mat.NewVecDense(2, []float64{1, 2})
	constr := mat.NewDense(1, 2, []float64{0, 0})
	rhs := mat.NewVecDense(1, []float64{0})
	scf := &StandardComputationalForm{
		Objective:      obj,
		Constraints:    constr,
		RHS:            rhs,
		PrimalSolution: mat.NewVecDense(2, nil),
	}

	scf.AddBranch(1, 5, 1)
	m, n := scf.Constraints.Dims()
	assert.Equal(t, m, 2)
	assert.Equal(t, n, 2)
	assert.Equal(t, scf.Constraints.At(1, 1), 1.0)
	assert.Equal(t, scf.RHS.AtVec(1), 5.0)

	scf.AddBranch(0, 3, 2)
	assert.Equal(t, scf.Constraints.At(2, 0), -1.0)
	assert.Equal(t, scf.RHS.AtVec(2), -3.0)
}

func TestSCFAddEquality(t *testing.T) {
	// create a minimal SCF
	scf := &StandardComputationalForm{
		Objective:      mat.NewVecDense(1, []float64{1}),
		Constraints:    mat.NewDense(1, 1, []float64{2}),
		RHS:            mat.NewVecDense(1, []float64{3}),
		PrimalSolution: mat.NewVecDense(1, []float64{0}),
		SlackIndices:   []int{-1},
		VarCategories:  []VarCategory{VarCategoryInteger},
	}
	scf.AddEquality(0, 5)
	if scf.Constraints.RawMatrix().Rows != 2 {
		t.Fatalf("expected 2 rows after AddEquality, got %d", scf.Constraints.RawMatrix().Rows)
	}
	if scf.RHS.Len() != 2 {
		t.Fatalf("expected RHS length 2 after AddEquality, got %d", scf.RHS.Len())
	}
	if scf.RHS.AtVec(1) != 5 {
		t.Fatalf("expected equality RHS 5, got %v", scf.RHS.AtVec(1))
	}
}
