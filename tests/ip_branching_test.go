package tests

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/internal/concurrency"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

// linExpr builds the expression sum(coefs[i] * vars[i]); variables without a
// coefficient are left out.
func linExpr(vars []lp.LpVariable, coefs ...float64) lp.LpExpression {
	terms := make([]lp.LpTerm, len(coefs))
	for i, c := range coefs {
		terms[i] = lp.NewTerm(c, vars[i])
	}
	return lp.NewExpression(terms)
}

// Branching on a general integer must bound the variable (x <= floor or
// x >= ceil) rather than fix it, otherwise optima away from the LP relaxation
// are cut off.
func Test_IPBranchingBoundsVariable(t *testing.T) {
	x := lp.NewVariable("x", lp.LpCategoryInteger)
	y := lp.NewVariable("y", lp.LpCategoryInteger)
	vars := []lp.LpVariable{x, y}

	t.Run("optimum below the floor of the relaxation", func(t *testing.T) {
		// LP relaxation: (2.25, 3.75). IP optimum: (0, 5) with objective 40.
		prog := lp.NewLinearProgram("below floor", vars)
		prog.AddObjective(lp.LpMaximise, linExpr(vars, 5, 8))
		prog.AddConstraint(linExpr(vars, 1, 1), lp.LpConstraintLE, 6)
		prog.AddConstraint(linExpr(vars, 5, 9), lp.LpConstraintLE, 45)

		sol, err := solver.Solve(&prog)
		assert.Nil(t, err)
		assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
		assert.IsClose(t, sol.ObjectiveValue, 40, 1e-5)
		assert.IsClose(t, sol.PrimalSolution.AtVec(0), 0, 1e-5)
		assert.IsClose(t, sol.PrimalSolution.AtVec(1), 5, 1e-5)
	})

	t.Run("optimum at a bound of the other variable", func(t *testing.T) {
		// LP relaxation: (3, 1.5). IP optimum: (4, 0) with objective 20.
		prog := lp.NewLinearProgram("other bound", vars)
		prog.AddObjective(lp.LpMaximise, linExpr(vars, 5, 4))
		prog.AddConstraint(linExpr(vars, 6, 4), lp.LpConstraintLE, 24)
		prog.AddConstraint(linExpr(vars, 1, 2), lp.LpConstraintLE, 6)

		sol, err := solver.Solve(&prog)
		assert.Nil(t, err)
		assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
		assert.IsClose(t, sol.ObjectiveValue, 20, 1e-5)
		assert.IsClose(t, sol.PrimalSolution.AtVec(0), 4, 1e-5)
		assert.IsClose(t, sol.PrimalSolution.AtVec(1), 0, 1e-5)
	})

	t.Run("up branch", func(t *testing.T) {
		// LP relaxation: x = 2.5. IP optimum: x = 3, reached via x >= 3.
		prog := lp.NewLinearProgram("up branch", vars)
		prog.AddObjective(lp.LpMinimise, linExpr(vars, 1, 0))
		prog.AddConstraint(linExpr(vars, 1, 0), lp.LpConstraintGE, 2.5)

		sol, err := solver.Solve(&prog)
		assert.Nil(t, err)
		assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
		assert.IsClose(t, sol.ObjectiveValue, 3, 1e-5)
		assert.IsClose(t, sol.PrimalSolution.AtVec(0), 3, 1e-5)
	})
}

// Simplex returns integral values a few ulps off (e.g. 0.99999999999999989).
// Integrality checks must allow for that instead of branching on them again,
// which adds a row per level and never terminates.
func Test_IPIntegralityTolerance(t *testing.T) {
	vars := []lp.LpVariable{
		lp.NewVariable("x0", lp.LpCategoryInteger),
		lp.NewVariable("x1", lp.LpCategoryInteger),
		lp.NewVariable("x2", lp.LpCategoryInteger),
	}
	prog := lp.NewLinearProgram("near integral", vars)
	prog.AddObjective(lp.LpMinimise, linExpr(vars, -4.1, 5.4, -2.9))
	prog.AddConstraint(linExpr(vars, -3.6, -8.3, 2.6), lp.LpConstraintLE, 14.8)
	prog.AddConstraint(linExpr(vars, 1, 0, 0), lp.LpConstraintLE, 7)
	prog.AddConstraint(linExpr(vars, 0, 1, 0), lp.LpConstraintLE, 7)
	prog.AddConstraint(linExpr(vars, 0, 0, 1), lp.LpConstraintLE, 7)

	type result struct {
		sol *solver.Solution
		err error
	}
	done := make(chan result, 1)
	go func() {
		sol, err := solver.Solve(&prog)
		done <- result{sol, err}
	}()

	select {
	case r := <-done:
		// Brute force over x_i in 0..7: optimum (7, 0, 7) with objective -49
		assert.Nil(t, r.err)
		assert.Equal(t, r.sol.Status.String(), lp.LpStatusOptimal.String())
		assert.IsClose(t, r.sol.ObjectiveValue, -49, 1e-5)
		assert.IsClose(t, r.sol.PrimalSolution.AtVec(0), 7, 1e-5)
		assert.IsClose(t, r.sol.PrimalSolution.AtVec(1), 0, 1e-5)
		assert.IsClose(t, r.sol.PrimalSolution.AtVec(2), 7, 1e-5)
	case <-time.After(10 * time.Second):
		t.Fatal("branch and bound did not terminate within 10s")
	}
}

// An integer solution is only accepted after an exact check against the
// original rows, so a relaxation value that is merely close to an integer
// cannot be reported as an optimum it does not satisfy.
func Test_IPExactAcceptance(t *testing.T) {
	t.Run("near-integral but infeasible", func(t *testing.T) {
		// LP relaxation: x = 0.99999999996666..., within the integrality tolerance
		// of 1, but x = 1 violates 3x <= 2.9999999999. The float simplex cannot
		// settle this (it also accepts the branch x >= 1), so the solve must fail
		// rather than report x = 1.
		vars := []lp.LpVariable{lp.NewVariable("x", lp.LpCategoryInteger)}
		prog := lp.NewLinearProgram("near integral", vars)
		prog.AddObjective(lp.LpMaximise, linExpr(vars, 1))
		prog.AddConstraint(linExpr(vars, 3), lp.LpConstraintLE, 2.9999999999)

		sol, err := solver.Solve(&prog)
		if err == nil {
			t.Fatalf("expected an error, got status=%v x=%v obj=%v", sol.Status, sol.PrimalSolution.AtVec(0), sol.ObjectiveValue)
		}
	})

	t.Run("decimal coefficients", func(t *testing.T) {
		// 0.1 + 0.2 <= 0.3 holds as written, although not for the float64 values.
		vars := []lp.LpVariable{
			lp.NewVariable("x", lp.LpCategoryInteger),
			lp.NewVariable("y", lp.LpCategoryInteger),
		}
		prog := lp.NewLinearProgram("decimal", vars)
		prog.AddObjective(lp.LpMaximise, linExpr(vars, 1, 1))
		prog.AddConstraint(linExpr(vars, 0.1, 0.2), lp.LpConstraintLE, 0.3)
		prog.AddConstraint(linExpr(vars, 1, 0), lp.LpConstraintLE, 1)
		prog.AddConstraint(linExpr(vars, 0, 1), lp.LpConstraintLE, 1)

		sol, err := solver.Solve(&prog)
		assert.Nil(t, err)
		assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
		assert.Equal(t, sol.ObjectiveValue, 2.0)
		assert.Equal(t, sol.PrimalSolution.AtVec(0), 1.0)
		assert.Equal(t, sol.PrimalSolution.AtVec(1), 1.0)
	})
}

// A child LP that fails (here: simplex iteration limit, as there is no
// anti-cycling rule yet) must fail the solve. Dropping its subtree silently
// reported this feasible IP as infeasible.
func Test_IPChildFailureIsNotSilent(t *testing.T) {
	vars := []lp.LpVariable{
		lp.NewVariable("x0", lp.LpCategoryInteger),
		lp.NewVariable("x1", lp.LpCategoryInteger),
		lp.NewVariable("x2", lp.LpCategoryInteger),
		lp.NewVariable("x3", lp.LpCategoryInteger),
	}
	prog := lp.NewLinearProgram("child failure", vars)
	prog.AddObjective(lp.LpMinimise, linExpr(vars, -3.8, -5.8, 5.3, -5.6))
	prog.AddConstraint(linExpr(vars, -9.1, -7.3, -8.4, 6), lp.LpConstraintGE, 3.6)
	prog.AddConstraint(linExpr(vars, 7.9, -8.7, -0.8, 3.1), lp.LpConstraintLE, 13.1)
	prog.AddConstraint(linExpr(vars, 5.5, 9.6, -9.3, 4.4), lp.LpConstraintLE, 4.4)
	for i := range vars {
		coefs := make([]float64, len(vars))
		coefs[i] = 1
		prog.AddConstraint(linExpr(vars, coefs...), lp.LpConstraintLE, 5)
	}

	sol, err := solver.Solve(&prog)
	if err != nil {
		t.Logf("solve failed loudly: %v", err)
		return
	}
	// Brute force over x_i in 0..5: optimum (0, 0, 3, 5) with objective -12.1
	assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
	assert.IsClose(t, sol.ObjectiveValue, -12.1, 1e-5)
}

// A better integer solution must win however close it is to the incumbent:
// neither the incumbent update nor pruning may give up optimality for a tolerance.
func Test_IPOptimalityWithoutTolerance(t *testing.T) {
	// Solve serially so the worse solution (x = 1, y = 0) is found first
	prev := atomic.LoadInt32(&concurrency.MAX_GOROUTINES)
	atomic.StoreInt32(&concurrency.MAX_GOROUTINES, 0)
	defer atomic.StoreInt32(&concurrency.MAX_GOROUTINES, prev)

	// max x + (1+5e-7) y  s.t. x + y <= 1.5, x <= 1, y <= 1
	// Integer points: x = 1, y = 0 with objective 1 and x = 0, y = 1 with 1.0000005.
	build := func(vars []lp.LpVariable, extra ...float64) *lp.LinearProgram {
		prog := lp.NewLinearProgram("close call", vars)
		prog.AddObjective(lp.LpMaximise, linExpr(vars, append([]float64{1, 1 + 5e-7}, extra...)...))
		prog.AddConstraint(linExpr(vars, 1, 1), lp.LpConstraintLE, 1.5)
		prog.AddConstraint(linExpr(vars, 1, 0), lp.LpConstraintLE, 1)
		prog.AddConstraint(linExpr(vars, 0, 1), lp.LpConstraintLE, 1)
		return &prog
	}
	x := lp.NewVariable("x", lp.LpCategoryInteger)
	y := lp.NewVariable("y", lp.LpCategoryInteger)

	t.Run("incumbent update", func(t *testing.T) {
		// x = 0, y = 1 is reached as an integer node after x = 1, y = 0
		sol, err := solver.Solve(build([]lp.LpVariable{x, y}))
		assert.Nil(t, err)
		assert.Equal(t, sol.ObjectiveValue, 1+5e-7)
		assert.Equal(t, sol.PrimalSolution.AtVec(1), 1.0)
	})

	t.Run("pruning", func(t *testing.T) {
		// With max ... + 1e-8 w and 2w <= 1, x = 0, y = 1 sits under the
		// fractional node w = 0.5, whose bound is less than 1e-6 above the incumbent
		vars := []lp.LpVariable{x, y, lp.NewVariable("w", lp.LpCategoryInteger)}
		prog := build(vars, 1e-8)
		prog.AddConstraint(linExpr(vars, 0, 0, 2), lp.LpConstraintLE, 1)
		sol, err := solver.Solve(prog)
		assert.Nil(t, err)
		assert.Equal(t, sol.ObjectiveValue, 1+5e-7)
		assert.Equal(t, sol.PrimalSolution.AtVec(1), 1.0)
	})
}
