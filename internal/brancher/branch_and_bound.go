package brancher

import (
	"fmt"

	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/internal/simplex"
)

func branchAndBound(ip *common.IntegerProgram, rootNode *common.Node, config *common.SolverConfig) error {
	nodes, err := branchFunc(rootNode)
	if err != nil {
		return errors.New(errors.ErrUnknown, "error in branching function", err)
	}

	for _, node := range nodes {
		node.Depth = rootNode.Depth + 1
		if config.Debug {
			fmt.Printf("[DEBUG] Branching to new node at depth %d\n", node.Depth)
		}
		err := simplex.Simplex(node.SCF, config)
		if err != nil {
			return errors.New(errors.ErrUnknown, "error solving child node", err)
		}

		if *node.SCF.Status != common.SolverStatusOptimal {
			// Node is infeasible, or unbounded, so it can be pruned
			continue
		}

		node.IsInteger = isIntegerFeasible(node.SCF)
		if config.Debug {
			fmt.Printf("[DEBUG] Node Objective: %.4f, IsInteger: %v\n\n", *node.SCF.ObjectiveValue, node.IsInteger)
			fmt.Printf("[DEBUG] Primal Solution: %v\n", node.SCF.PrimalSolution)
		}
		if node.IsInteger {
			objVal := *node.SCF.ObjectiveValue
			if objVal < ip.BestObj+config.Tolerance {
				ip.BestObj = objVal
				ip.BestSolution = node.SCF.PrimalSolution
				if config.Debug {
					fmt.Printf("[DEBUG] New Best Obj: %.4f\n", ip.BestObj)
				}
			}
			continue
		}

		// If not integer feasible, continue branching
		err = branchAndBound(ip, node, config)
		if err != nil && config.Logging {
			fmt.Printf("Error in branchAndBound: %v\n", err)
		}
	}

	return nil
}
