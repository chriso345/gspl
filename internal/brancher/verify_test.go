package brancher

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/internal/common"
	"gonum.org/v1/gonum/mat"
)

func TestVerifyIntegerSolution(t *testing.T) {
	integer := common.VarCategoryInteger
	continuous := common.VarCategoryContinuous

	t.Run("violated by less than the tolerance", func(t *testing.T) {
		// 3x + s = 2.9999999999: x = 1 violates it by 1e-10
		scf := &common.StandardComputationalForm{
			Objective:     mat.NewVecDense(2, []float64{-1, 0}),
			Constraints:   mat.NewDense(1, 2, []float64{3, 1}),
			RHS:           mat.NewVecDense(1, []float64{2.9999999999}),
			SlackIndices:  []int{-1, 1},
			VarCategories: []common.VarCategory{integer, continuous},
		}
		_, _, ok := verifyIntegerSolution(scf, mat.NewVecDense(2, []float64{0.99999999996666666, 0}))
		assert.False(t, ok)

		point, obj, ok := verifyIntegerSolution(scf, mat.NewVecDense(2, []float64{0, 2.9999999999}))
		assert.True(t, ok)
		assert.Equal(t, point.AtVec(0), 0.0)
		assert.Equal(t, obj, 0.0)
	})

	t.Run("decimal coefficients", func(t *testing.T) {
		// 0.1x + 0.2y + s = 0.3 holds for x = y = 1 as written, not for the float64 values
		scf := &common.StandardComputationalForm{
			Objective:     mat.NewVecDense(3, []float64{-1, -1, 0}),
			Constraints:   mat.NewDense(1, 3, []float64{0.1, 0.2, 1}),
			RHS:           mat.NewVecDense(1, []float64{0.3}),
			SlackIndices:  []int{-1, -1, 2},
			VarCategories: []common.VarCategory{integer, integer, continuous},
		}
		point, obj, ok := verifyIntegerSolution(scf, mat.NewVecDense(3, []float64{1, 0.99999999999999989, 0}))
		assert.True(t, ok)
		assert.Equal(t, point.AtVec(1), 1.0)
		assert.Equal(t, obj, -2.0)
	})

	t.Run("surplus and equality rows", func(t *testing.T) {
		// x - s = 2 (x >= 2) and x + y = 3
		scf := &common.StandardComputationalForm{
			Objective:     mat.NewVecDense(3, []float64{0, 0, 0}),
			Constraints:   mat.NewDense(2, 3, []float64{1, 0, -1, 1, 1, 0}),
			RHS:           mat.NewVecDense(2, []float64{2, 3}),
			SlackIndices:  []int{-1, -1, 2},
			VarCategories: []common.VarCategory{integer, integer, continuous},
		}
		_, _, ok := verifyIntegerSolution(scf, mat.NewVecDense(3, []float64{2, 1, 0}))
		assert.True(t, ok)
		_, _, ok = verifyIntegerSolution(scf, mat.NewVecDense(3, []float64{1, 2, 0}))
		assert.False(t, ok)
		_, _, ok = verifyIntegerSolution(scf, mat.NewVecDense(3, []float64{2, 2, 0}))
		assert.False(t, ok)
	})

	t.Run("binary out of range", func(t *testing.T) {
		scf := &common.StandardComputationalForm{
			Objective:     mat.NewVecDense(2, []float64{-1, 0}),
			Constraints:   mat.NewDense(1, 2, []float64{1, 1}),
			RHS:           mat.NewVecDense(1, []float64{5}),
			SlackIndices:  []int{-1, 1},
			VarCategories: []common.VarCategory{common.VarCategoryBinary, continuous},
		}
		_, _, ok := verifyIntegerSolution(scf, mat.NewVecDense(2, []float64{2, 3}))
		assert.False(t, ok)
	})
}
