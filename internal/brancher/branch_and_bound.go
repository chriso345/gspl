package brancher

import (
	"fmt"

	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/concurrency"
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/internal/simplex"
)

func branchAndBound(ip *common.IntegerProgram, rootNode *common.Node, config *common.SolverConfig) error {
	return branchAndBoundParallel(ip, rootNode, config)
}

// branchAndBoundParallel runs branch-and-bound in parallel using goroutines and channels.
func branchAndBoundParallel(ip *common.IntegerProgram, rootNode *common.Node, config *common.SolverConfig) error {
	nodes, err := branchFunc(rootNode)
	if err != nil {
		return errors.New(errors.ErrUnknown, "error in branching function", err)
	}

	type result struct {
		node *common.Node
		err  error
	}
	results := make(chan result, len(nodes))

	// helper to process a node (can run inline or in goroutine)
	processNode := func(root *common.Node, node *common.Node) error {
		node.Depth = root.Depth + 1
		if config.Debug {
			fmt.Printf("[DEBUG] Branching to new node at depth %d\n", node.Depth)
		}
		err := simplex.Simplex(node.SCF, config)
		if err != nil {
			return err
		}
		if *node.SCF.Status != common.SolverStatusOptimal {
			return nil
		}
		node.IsInteger = isIntegerFeasible(node.SCF)
		if config.Debug {
			fmt.Printf("[DEBUG] Node Objective: %.4f, IsInteger: %v\n\n", *node.SCF.ObjectiveValue, node.IsInteger)
			fmt.Printf("[DEBUG] Primal Solution: %v\n", node.SCF.PrimalSolution)
		}
		if node.IsInteger {
			objVal := *node.SCF.ObjectiveValue
			// Flip to original sense if this SCF represents a maximisation
			if node.SCF.IsMaximization {
				objVal = -objVal
			}
			// protect BestObj update
			ip.BestMutex.Lock()
			// If no best solution yet, accept this one
			if ip.BestSolution == nil {
				ip.BestObj = objVal
				ip.BestSolution = node.SCF.PrimalSolution
				if config.Debug {
					fmt.Printf("[DEBUG] New Best Obj: %.4f\n", ip.BestObj)
				}
				ip.BestMutex.Unlock()
				return nil
			}
			// Update depending on minimisation/maximisation
			if node.SCF.IsMaximization {
				if objVal > ip.BestObj+config.Tolerance {
					ip.BestObj = objVal
					ip.BestSolution = node.SCF.PrimalSolution
					if config.Debug {
						fmt.Printf("[DEBUG] New Best Obj: %.4f\n", ip.BestObj)
					}
				}
			} else {
				if objVal < ip.BestObj-config.Tolerance {
					ip.BestObj = objVal
					ip.BestSolution = node.SCF.PrimalSolution
					if config.Debug {
						fmt.Printf("[DEBUG] New Best Obj: %.4f\n", ip.BestObj)
					}
				}
			}
			ip.BestMutex.Unlock()
			return nil
		}
		// Not integer feasible, branch recursively
		return branchAndBoundParallel(ip, node, config)
	}

	for _, node := range nodes {
		nd := node
		// try to acquire goroutine slot
		if concurrency.TryAcquireGoroutine() {
			go func(n *common.Node) {
				defer concurrency.ReleaseGoroutine()
				err := processNode(rootNode, n)
				results <- result{n, err}
			}(nd)
		} else {
			// run inline
			err := processNode(rootNode, nd)
			results <- result{nd, err}
		}
	}

	for range nodes {
		r := <-results
		if r.err != nil && config.Logging {
			fmt.Printf("Error in branchAndBoundParallel: %v\n", r.err)
		}
	}
	return nil
}
