package tests

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

func Test_IPMinimisationExample(t *testing.T) {
	variables := []lp.LpVariable{
		lp.NewVariable("x1", lp.LpCategoryInteger),
		lp.NewVariable("x2", lp.LpCategoryInteger),
	}

	// Objective: Minimize 3x1 + 2x2
	objTerms := []lp.LpTerm{
		lp.NewTerm(3, variables[0]),
		lp.NewTerm(2, variables[1]),
	}
	objective := lp.NewExpression(objTerms)

	// Constraints:
	// 1.5x1 + x2 >= 7
	con1Terms := []lp.LpTerm{
		lp.NewTerm(1.5, variables[0]),
		lp.NewTerm(1, variables[1]),
	}

	// x1 + 0.5x2 >= 3
	con2Terms := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(0.5, variables[1]),
	}

	// Build LP minimization problem
	lpProg := lp.NewLinearProgram("Non-Integer LP Example", variables)
	lpProg.AddObjective(lp.LpMinimise, objective)
	lpProg.AddConstraint(lp.NewExpression(con1Terms), lp.LpConstraintGE, 7)
	lpProg.AddConstraint(lp.NewExpression(con2Terms), lp.LpConstraintGE, 3)

	// Solve it
	sol, err := solver.Solve(&lpProg)
	assert.Nil(t, err)
	assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
	assert.IsClose(t, sol.ObjectiveValue, 14.0, 1e-5)
}
