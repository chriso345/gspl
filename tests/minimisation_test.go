package tests

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

func Test_MinimisationExample(t *testing.T) {
	variables := []lp.LpVariable{
		lp.NewVariable("x1"),
		lp.NewVariable("x2"),
		lp.NewVariable("x3"),
		lp.NewVariable("x4"),
		lp.NewVariable("x5"),
	}

	terms := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(2, variables[1]),
		lp.NewTerm(3, variables[2]),
		lp.NewTerm(1, variables[3]),
		lp.NewTerm(4, variables[4]),
	}
	objective := lp.NewExpression(terms)

	terms2 := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(1, variables[1]),
		lp.NewTerm(1, variables[2]),
		lp.NewTerm(1, variables[3]),
		lp.NewTerm(1, variables[4]),
	}

	terms3 := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(2, variables[1]),
		lp.NewTerm(1, variables[2]),
		lp.NewTerm(0, variables[3]),
		lp.NewTerm(0, variables[4]),
	}

	terms4 := []lp.LpTerm{
		lp.NewTerm(0, variables[0]),
		lp.NewTerm(1, variables[1]),
		lp.NewTerm(0, variables[2]),
		lp.NewTerm(1, variables[3]),
		lp.NewTerm(1, variables[4]),
	}

	terms5 := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(0, variables[1]),
		lp.NewTerm(1, variables[2]),
		lp.NewTerm(0, variables[3]),
		lp.NewTerm(1, variables[4]),
	}

	terms6 := []lp.LpTerm{
		lp.NewTerm(1, variables[3]),
		lp.NewTerm(1, variables[4]),
	}

	minProg := lp.NewLinearProgram("Minimisation Example", variables)
	minProg.AddObjective(lp.LpMinimise, objective)
	minProg.AddConstraint(lp.NewExpression(terms2), lp.LpConstraintGE, 10)
	minProg.AddConstraint(lp.NewExpression(terms3), lp.LpConstraintLE, 8)
	minProg.AddConstraint(lp.NewExpression(terms4), lp.LpConstraintLE, 7)
	minProg.AddConstraint(lp.NewExpression(terms5), lp.LpConstraintGE, 4)
	minProg.AddConstraint(lp.NewExpression(terms6), lp.LpConstraintLE, 6)

	sol, err := solver.Solve(&minProg)

	assert.Nil(t, err)
	assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
	assert.IsClose(t, sol.ObjectiveValue, 10.0, 1e-5)
}

func Test_Minimisation1(t *testing.T) {
	variables := []lp.LpVariable{
		lp.NewVariable("x1", lp.LpCategoryInteger),
		lp.NewVariable("x2", lp.LpCategoryInteger),
		lp.NewVariable("x3", lp.LpCategoryInteger),
	}

	objTerms := []lp.LpTerm{
		lp.NewTerm(2, variables[0]),
		lp.NewTerm(3, variables[1]),
		lp.NewTerm(4, variables[2]),
	}
	objective := lp.NewExpression(objTerms)

	con1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[0]), lp.NewTerm(1, variables[1])})
	con2 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[1]), lp.NewTerm(1, variables[2])})
	con3 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[0]), lp.NewTerm(1, variables[2])})
	con4 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[0]), lp.NewTerm(1, variables[1]), lp.NewTerm(1, variables[2])})
	con5 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(2, variables[0]), lp.NewTerm(1, variables[1]), lp.NewTerm(1, variables[2])})

	prog := lp.NewLinearProgram("Medium Benchmark", variables)
	prog.AddObjective(lp.LpMinimise, objective)

	prog.AddConstraint(con1, lp.LpConstraintGE, 1)
	prog.AddConstraint(con2, lp.LpConstraintGE, 1)
	prog.AddConstraint(con3, lp.LpConstraintGE, 1)
	prog.AddConstraint(con4, lp.LpConstraintGE, 2)
	prog.AddConstraint(con5, lp.LpConstraintGE, 3)

	sol, err := solver.Solve(&prog)

	assert.Nil(t, err)
	assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
	assert.IsClose(t, sol.ObjectiveValue, 5.0, 1e-5)
}

func Test_Minimisation2(t *testing.T) {
	variables := []lp.LpVariable{
		lp.NewVariable("x1", lp.LpCategoryInteger),
		lp.NewVariable("x2", lp.LpCategoryInteger),
		lp.NewVariable("x3", lp.LpCategoryInteger),
		lp.NewVariable("x4", lp.LpCategoryInteger),
		lp.NewVariable("x5", lp.LpCategoryInteger),
	}

	objTerms := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(2, variables[1]),
		lp.NewTerm(3, variables[2]),
		lp.NewTerm(4, variables[3]),
		lp.NewTerm(5, variables[4]),
	}
	objective := lp.NewExpression(objTerms)

	constraints := []lp.LpExpression{
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[0]), lp.NewTerm(1, variables[1])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[1]), lp.NewTerm(1, variables[2])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[2]), lp.NewTerm(1, variables[3])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[3]), lp.NewTerm(1, variables[4])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[0]), lp.NewTerm(1, variables[4])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[0]), lp.NewTerm(1, variables[2])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[1]), lp.NewTerm(1, variables[3])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[2]), lp.NewTerm(1, variables[4])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[0]), lp.NewTerm(1, variables[1]), lp.NewTerm(1, variables[2])}),
		lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, variables[1]), lp.NewTerm(1, variables[3]), lp.NewTerm(1, variables[4])}),
	}

	constValues := []float64{1, 1, 1, 1, 1, 2, 2, 2, 2, 2}

	prog := lp.NewLinearProgram("Large Benchmark", variables)
	prog.AddObjective(lp.LpMinimise, objective)

	for i := range 10 {
		prog.AddConstraint(constraints[i], lp.LpConstraintGE, constValues[i])
	}

	sol, err := solver.Solve(&prog)

	assert.Nil(t, err)
	assert.Equal(t, sol.Status.String(), lp.LpStatusOptimal.String())
	assert.IsClose(t, sol.ObjectiveValue, 10.0, 1e-5)
}
