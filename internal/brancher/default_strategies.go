package brancher

import (
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
)

// DefaultBranch represents the default branching strategy.
//
// This branches on the first variable found that is not integer in the current node
func DefaultBranch(node *common.Node) ([]*common.Node, error) {
	// Build mapping from primal indices to original variable indices
	primalToVar := []int{}
	for varIdx, slackMarker := range node.SCF.SlackIndices {
		if slackMarker == -1 {
			primalToVar = append(primalToVar, varIdx)
		}
	}

	// Find branching variable (first non-integer among primals or binary out-of-range)
	branchingVarIndex := -1
	for p := 0; p < node.SCF.PrimalSolution.Len(); p++ {
		val := node.SCF.PrimalSolution.AtVec(p)
		// Map primal index to variable index
		var varIdx int
		if p < len(primalToVar) {
			varIdx = primalToVar[p]
		} else {
			varIdx = p
		}
		isBinary := varIdx < len(node.SCF.VarCategories) && node.SCF.VarCategories[varIdx] == common.VarCategoryBinary
		if val != float64(int(val)) || (isBinary && !(val == 0 || val == 1)) {
			branchingVarIndex = p
			break
		}
	}

	if branchingVarIndex == -1 {
		return nil, errors.New(errors.ErrInfeasible, "no branching variable found; node is already integer feasible", nil)
	}

	val := node.SCF.PrimalSolution.AtVec(branchingVarIndex)
	varIndex := branchingVarIndex
	if branchingVarIndex < len(primalToVar) {
		varIndex = primalToVar[branchingVarIndex]
	}

	// If variable category is continuous, skip branching on it (MILP allows continuous non-integers)
	if varIndex < len(node.SCF.VarCategories) && node.SCF.VarCategories[varIndex] == common.VarCategoryContinuous {
		// Do not branch on continuous variables, find next
		found := false
		for p := branchingVarIndex + 1; p < node.SCF.PrimalSolution.Len(); p++ {
			val2 := node.SCF.PrimalSolution.AtVec(p)
			if val2 != float64(int(val2)) {
				vIdx := p
				if p < len(primalToVar) {
					vIdx = primalToVar[p]
				}
				if p < len(node.SCF.VarCategories) && node.SCF.VarCategories[vIdx] != common.VarCategoryContinuous {
					branchingVarIndex = p
					val = val2
					varIndex = vIdx
					found = true
					break
				}
			}
		}
		if !found {
			return nil, errors.New(errors.ErrInfeasible, "no suitable branching variable found; node may be MILP-feasible", nil)
		}
	}

	// If this column is binary, create equality fixes to 0 and 1.
	if varIndex < len(node.SCF.VarCategories) && node.SCF.VarCategories[varIndex] == common.VarCategoryBinary {
		down := &common.Node{SCF: node.SCF.Copy()}
		up := &common.Node{SCF: node.SCF.Copy()}
		// Fix to 0 and 1; AddEquality expects column index in SCF (primal index corresponds to column)
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
