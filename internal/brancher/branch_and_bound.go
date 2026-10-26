package brancher

import (
	"fmt"
	"math"

	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/concurrency"
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/internal/simplex"
)

// pruneMargin is how far, relative to the incumbent, a node's LP bound must be
// below it (above it when minimising) before the node is pruned. It covers float
// error in the bound, so it errs towards exploring a node.
const pruneMargin = 1e-9

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
			return acceptIncumbent(ip, node.SCF, config)
		}
		// Not integer feasible: prune if its LP bound is clearly worse than the
		// incumbent (see pruneMargin), otherwise branch recursively. The bound is
		// in the SCF's minimisation sense.
		ip.BestMutex.Lock()
		best := ip.BestObj
		if node.SCF.IsMaximization {
			best = -best
		}
		dominated := ip.BestSolution != nil &&
			*node.SCF.ObjectiveValue > best+pruneMargin*math.Max(1, math.Abs(best))
		ip.BestMutex.Unlock()
		if dominated {
			return nil
		}
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

	// Wait for every child, then report the first failure: a dropped subtree
	// could hide the optimum.
	var firstErr error
	for range nodes {
		r := <-results
		if r.err != nil {
			if config.Logging {
				fmt.Printf("Error in branchAndBoundParallel: %v\n", r.err)
			}
			if firstErr == nil {
				firstErr = r.err
			}
		}
	}
	return firstErr
}
