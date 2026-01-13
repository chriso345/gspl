package brancher

import (
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
)

// DefaultBranch represents the default branching strategy.
//
// This branches on the first variable found that is not integer in the current node
func DefaultBranch(node *common.Node) ([]*common.Node, error) {
	// Find branching variable
	branchingVarIndex := -1
	for i := 0; i < node.SCF.PrimalSolution.Len(); i++ {
		val := node.SCF.PrimalSolution.AtVec(i)
		if val != float64(int(val)) {
			branchingVarIndex = i
			break
		}
	}

	if branchingVarIndex == -1 {
		return nil, errors.New(errors.ErrInfeasible, "no branching variable found; node is already integer feasible", nil)
	}

	val := node.SCF.PrimalSolution.AtVec(branchingVarIndex)

	// If variable category is continuous, skip branching on it (MILP allows continuous non-integers)
	if branchingVarIndex < len(node.SCF.VarCategories) && node.SCF.VarCategories[branchingVarIndex] == common.VarCategoryContinuous {
		// Do not branch on continuous variables, find next
		found := false
		for i := branchingVarIndex + 1; i < node.SCF.PrimalSolution.Len(); i++ {
			val2 := node.SCF.PrimalSolution.AtVec(i)
			if val2 != float64(int(val2)) && i < len(node.SCF.VarCategories) && node.SCF.VarCategories[i] != common.VarCategoryContinuous {
				branchingVarIndex = i
				val = val2
				found = true
				break
			}
		}
		if !found {
			return nil, errors.New(errors.ErrInfeasible, "no suitable branching variable found; node may be MILP-feasible", nil)
		}
	}

	// If this column is binary, create equality fixes to 0 and 1.
	if branchingVarIndex < len(node.SCF.VarCategories) && node.SCF.VarCategories[branchingVarIndex] == common.VarCategoryBinary {
		down := &common.Node{SCF: node.SCF.Copy()}
		up := &common.Node{SCF: node.SCF.Copy()}
		// Fix to 0 and 1
		down.SCF.AddEquality(branchingVarIndex, 0)
		up.SCF.AddEquality(branchingVarIndex, 1)
		return []*common.Node{up, down}, nil
	}

	// Fallback: general integer branching using AddBranch
	down := &common.Node{SCF: node.SCF.Copy()}
	up := &common.Node{SCF: node.SCF.Copy()}
	down.SCF.AddBranch(branchingVarIndex, float64(int(val)), 1)
	up.SCF.AddBranch(branchingVarIndex, float64(int(val)+1), 2)
	return []*common.Node{up, down}, nil
}

// DefaultHeuristic represents the default heuristic strategy.
//
// This does not implement any heuristic and simply returns nil
func DefaultHeuristic(node *common.Node) ([]float64, float64, bool) {
	return nil, 0, false
}

// DefaultCut represents the default cutting planes strategy.
//
// This does not implement any cutting planes and simply returns nil
func DefaultCut(node *common.Node) [][]float64 {
	return nil
}
