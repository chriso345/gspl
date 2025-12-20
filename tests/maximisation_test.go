package tests

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

func Test_MaximisationExample(t *testing.T) {
	// Continuous decision variables
	variables := []lp.LpVariable{
		lp.NewVariable("x1", lp.LpCategoryContinuous),
		lp.NewVariable("x2", lp.LpCategoryContinuous),
	}

	// Objective: Maximize 5*x1 + 4*x2
	objTerms := []lp.LpTerm{
		lp.NewTerm(5, variables[0]),
		lp.NewTerm(4, variables[1]),
	}
	objective := lp.NewExpression(objTerms)

	// Constraints:
	// 2*x1 + 3*x2 <= 12
	con1Terms := []lp.LpTerm{
		lp.NewTerm(2, variables[0]),
		lp.NewTerm(3, variables[1]),
	}
	// x1 + x2 <= 5
	con2Terms := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(1, variables[1]),
	}

	// Build LP maximization problem
	lpProg := lp.NewLinearProgram("Maximization Example", variables)
	lpProg.AddObjective(lp.LpMaximise, objective)
	lpProg.AddConstraint(lp.NewExpression(con1Terms), lp.LpConstraintLE, 12)
	lpProg.AddConstraint(lp.NewExpression(con2Terms), lp.LpConstraintLE, 5)

	// Solve it
	sol, err := solver.Solve(&lpProg)
	assert.Nil(t, err)
	assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
	assert.IsClose(t, sol.ObjectiveValue, 25.0, 1e-5)
}

func Test_SimpleMaximisation(t *testing.T) {
	// --- Decision variables ---
	variables := []lp.LpVariable{
		lp.NewVariable("y1"),
		lp.NewVariable("y2"),
		lp.NewVariable("y3"),
		lp.NewVariable("y4"),
	}

	// --- Objective: maximize -9*y1 - 18*y2 - 7*y3 - 6*y4 ---
	objTerms := []lp.LpTerm{
		lp.NewTerm(-9, variables[0]),
		lp.NewTerm(-18, variables[1]),
		lp.NewTerm(-7, variables[2]),
		lp.NewTerm(-6, variables[3]),
	}
	objective := lp.NewExpression(objTerms)

	// --- Constraints ---
	con1 := lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(3, variables[1]),
		lp.NewTerm(1, variables[2]),
	})
	con2 := lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(1, variables[1]),
		lp.NewTerm(1, variables[3]),
	})

	// --- Build Linear Program ---
	prog := lp.NewLinearProgram("Simple Maximisation", variables)
	prog.AddObjective(lp.LpMaximise, objective)
	prog.AddConstraint(con1, lp.LpConstraintGE, 3)
	prog.AddConstraint(con2, lp.LpConstraintGE, 2)

	// --- Solve ---
	sol, err := solver.Solve(&prog)

	// --- Assertions ---
	assert.Nil(t, err)
	assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
	assert.IsClose(t, sol.ObjectiveValue, -22.5, 1e-5)
}
