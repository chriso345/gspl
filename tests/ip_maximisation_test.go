package tests

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

func Test_KnapsackProblem(t *testing.T) {
	variables := []lp.LpVariable{
		lp.NewVariable("x1", lp.LpCategoryInteger),
		lp.NewVariable("x2", lp.LpCategoryInteger),
		lp.NewVariable("x3", lp.LpCategoryInteger),
		lp.NewVariable("x4", lp.LpCategoryInteger),
		lp.NewVariable("x5", lp.LpCategoryInteger),
	}

	// --- Objective ---
	objTerms := []lp.LpTerm{
		lp.NewTerm(5, variables[0]),
		lp.NewTerm(3, variables[1]),
		lp.NewTerm(6, variables[2]),
		lp.NewTerm(6, variables[3]),
		lp.NewTerm(2, variables[4]),
	}
	objective := lp.NewExpression(objTerms)

	// --- Knapsack constraint ---
	knapsackTerms := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(4, variables[1]),
		lp.NewTerm(7, variables[2]),
		lp.NewTerm(6, variables[3]),
		lp.NewTerm(2, variables[4]),
	}
	knapsackConstraint := lp.NewExpression(knapsackTerms)

	prog := lp.NewLinearProgram("Knapsack Problem", variables)
	prog.AddObjective(lp.LpMaximise, objective)
	prog.AddConstraint(knapsackConstraint, lp.LpConstraintLE, 15)

	// --- Variable bounds as constraints ---
	for _, v := range variables {
		// x_i >= 0
		prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, v)}), lp.LpConstraintGE, 0)
		// x_i <= 1
		prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, v)}), lp.LpConstraintLE, 1)
	}

	sol, err := solver.Solve(&prog)

	assert.Nil(t, err)
	assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
	assert.IsClose(t, sol.ObjectiveValue, 17, 1e-5)
}
