package brancher

import (
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
)

// DefaultBranch represents the default branching strategy.
//
// This branches on the first variable found that is not integer in the current node
func DefaultBranch(node *common.Node) ([]*common.Node, error) {
	down := &common.Node{
		SCF: node.SCF.Copy(),
	}
	up := &common.Node{
		SCF: node.SCF.Copy(),
	}

	// Get the branching variable
	branchingVarIndex := -1
	for i := 0; i < node.SCF.PrimalSolution.Len(); i++ {
		val := node.SCF.PrimalSolution.AtVec(i)
		if val != float64(int(val)) {
			branchingVarIndex = i
			break
		}
	}

	// Add constraints to respective child nodes (added as <= constraints)
	down.SCF.AddBranch(branchingVarIndex, float64(int(node.SCF.PrimalSolution.AtVec(branchingVarIndex))), 1)
	up.SCF.AddBranch(branchingVarIndex, float64(int(node.SCF.PrimalSolution.AtVec(branchingVarIndex))+1), 2)

	if branchingVarIndex == -1 {
		return nil, errors.New(errors.ErrInfeasible, "no branching variable found; node is already integer feasible", nil)
	}

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
