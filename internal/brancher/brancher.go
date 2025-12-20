package brancher

import (
	"fmt"
	"math"

	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/internal/simplex"
)

func BranchAndBound(ip *common.IntegerProgram, config *common.SolverConfig) error {
	// Define the strategies to be used in tree traversal
	defineStrategies(ip)

	// Solve at the root
	rootNode := &common.Node{
		SCF: ip.SCF,
		// ID:       0,
		// ParentID: -1,
		// Depth:    0,
	}

	err := simplex.Simplex(rootNode.SCF, config)
	if err != nil {
		return errors.New(errors.ErrUnknown, "error solving root node", err)
	}

	// If the root node is not optimal, the IP is infeasible or unbounded
	if *rootNode.SCF.Status != common.SolverStatusOptimal {
		*ip.SCF.Status = *rootNode.SCF.Status
		return nil
	}

	// Check if the root solution is integer feasible
	rootNode.IsInteger = isIntegerFeasible(rootNode.SCF)

	if config.Logging {
		fmt.Printf("[DEBUG] Primal Solution: %v\n", rootNode.SCF.PrimalSolution)
	}

	if rootNode.IsInteger {
		// BestObj is stored in the original problem sense. If the SCF indicates
		// the original problem was a maximisation, flip the sign (Simplex returns
		// its minimisation-form objective value).
		if rootNode.SCF.IsMaximization {
			ip.BestObj = -(*rootNode.SCF.ObjectiveValue)
		} else {
			ip.BestObj = *rootNode.SCF.ObjectiveValue
		}
		ip.BestSolution = rootNode.SCF.PrimalSolution
		*ip.SCF.Status = common.SolverStatusOptimal
		return nil
	}

	err = branchAndBound(ip, rootNode, config)

	// Set final SCF status depending on whether a best solution was found
	if ip.BestSolution != nil {
		*ip.SCF.Status = common.SolverStatusOptimal
	} else {
		*ip.SCF.Status = common.SolverStatusInfeasible
	}

	ip.SCF.ObjectiveValue = &ip.BestObj
	ip.SCF.PrimalSolution = ip.BestSolution
	return err
}

// isIntegerFeasible checks if a solution is currently integer feasible
func isIntegerFeasible(scf *common.StandardComputationalForm) bool {
	sol := scf.PrimalSolution
	for i := 0; i < sol.Len(); i++ {
		val := sol.AtVec(i)
		// If the current value is not integer and it the index is not in the slack indices
		isSlack := scf.SlackIndices[i]

		if math.Floor(val) != val && isSlack == -1 {
			return false
		}
	}
	return true
}

// defineStrategies sets the strategies to be used in the Branch and Bound algorithm
func defineStrategies(ip *common.IntegerProgram) {
	if ip.Branch == nil {
		branchFunc = DefaultBranch
	} else {
		branchFunc = ip.Branch
	}

	// Heuristic and Cut functions are read directly from the IntegerProgram where required.
	// If not provided, callers should use DefaultHeuristic/DefaultCut explicitly.
}

// Strategy function variables
var branchFunc common.BranchFunc
