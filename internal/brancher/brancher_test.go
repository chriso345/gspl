package brancher

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/internal/common"
	"gonum.org/v1/gonum/mat"
)

func TestIsIntegerFeasible(t *testing.T) {
	scf := newTestSCF([]float64{1, 2, 3})
	for i := range scf.SlackIndices {
		scf.SlackIndices[i] = -1
	}
	assert.True(t, isIntegerFeasible(scf))

	scf.PrimalSolution.SetVec(1, 2.5)
	assert.False(t, isIntegerFeasible(scf))
}

func TestDefaultBranch_BinaryCreatesEquality(t *testing.T) {
	scf := &common.StandardComputationalForm{
		PrimalSolution: mat.NewVecDense(2, []float64{0.2, 1.0}),
		VarCategories:  []common.VarCategory{common.VarCategoryBinary, common.VarCategoryInteger},
		SlackIndices:   []int{-1, -1},
		Constraints:    mat.NewDense(1, 2, []float64{0, 0}),
		RHS:            mat.NewVecDense(1, []float64{0}),
		Objective:      mat.NewVecDense(2, []float64{0, 0}),
		ObjectiveValue: new(float64),
		Status:         new(common.SolverStatus),
	}
	node := &common.Node{SCF: scf}
	children, err := DefaultBranch(node)
	assert.Nil(t, err)
	assert.Equal(t, len(children), 2)
	// children SCFs should have extra rows added by AddEquality
	assert.True(t, children[0].SCF.RHS.Len() >= 2)
}

func TestIsIntegerFeasible_VarCategoriesFallbackBinary(t *testing.T) {
	// No SlackIndices populated -> triggers fallback that uses VarCategories
	scf := &common.StandardComputationalForm{
		PrimalSolution: mat.NewVecDense(1, []float64{0.5}),
		SlackIndices:   []int{},
		VarCategories:  []common.VarCategory{common.VarCategoryBinary},
	}
	assert.False(t, isIntegerFeasible(scf))
}

func TestIsIntegerFeasible_MappingAndBinaryChecks(t *testing.T) {
	scf := &common.StandardComputationalForm{
		PrimalSolution: mat.NewVecDense(3, []float64{1.0, 1.0, 2.0}),
		SlackIndices:   []int{-1, -1, -1},
		VarCategories:  []common.VarCategory{common.VarCategoryInteger, common.VarCategoryBinary, common.VarCategoryInteger},
	}
	assert.True(t, isIntegerFeasible(scf))

	scf.PrimalSolution.SetVec(1, 0.5)
	assert.False(t, isIntegerFeasible(scf))
}

func TestIsIntegerFeasible_FallbackIntegerNonInteger(t *testing.T) {
	scf := &common.StandardComputationalForm{
		PrimalSolution: mat.NewVecDense(2, []float64{1.0, 2.5}),
		SlackIndices:   []int{0, 0}, // no -1 entries -> primalToVar len == 0 triggers fallback
		VarCategories:  []common.VarCategory{common.VarCategoryInteger, common.VarCategoryInteger},
	}
	assert.False(t, isIntegerFeasible(scf))
}
