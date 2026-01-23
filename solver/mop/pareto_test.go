package mop

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/lp"
)

func Test_ParetoWeightedSum2Objectives(t *testing.T) {
	x := lp.NewVariable("x", lp.LpCategoryContinuous)
	y := lp.NewVariable("y", lp.LpCategoryContinuous)
	variables := []lp.LpVariable{x, y}

	// Objective 1: maximize x
	obj1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)})
	// Objective 2: maximize y
	obj2 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, y)})

	con := lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(1, x),
		lp.NewTerm(1, y),
	})

	prog := lp.NewLinearProgram("Pareto Weighted Sum Test", variables)
	prog.AddObjective(lp.LpMaximise, obj1)
	prog.AddObjective(lp.LpMaximise, obj2)
	prog.AddConstraint(con, lp.LpConstraintLE, 10)

	// Solve using weighted sum
	sols, err := SolvePareto(&prog, ParetoWeightedSum)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(sols), 2)

	// Each solution should satisfy x + y = 10
	for _, s := range sols {
		sum := s.ObjectiveValues[0] + s.ObjectiveValues[1]
		assert.IsClose(t, sum, 10.0, 1e-5)
	}
}

func Test_ParetoEpsilonConstraint2Objectives(t *testing.T) {
	x := lp.NewVariable("x", lp.LpCategoryContinuous)
	y := lp.NewVariable("y", lp.LpCategoryContinuous)
	variables := []lp.LpVariable{x, y}

	// Objectives
	obj1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)})
	obj2 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, y)})

	con := lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(1, x),
		lp.NewTerm(1, y),
	})

	prog := lp.NewLinearProgram("Pareto Epsilon Constraint Test", variables)
	prog.AddObjective(lp.LpMaximise, obj1)
	prog.AddObjective(lp.LpMaximise, obj2)
	prog.AddConstraint(con, lp.LpConstraintLE, 10)

	// Solve using epsilon constraint method with 5 steps
	sols, err := SolvePareto(
		&prog,
		ParetoEpsilonConstraint,
	)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(sols), 2)

	for _, s := range sols {
		sum := s.ObjectiveValues[0] + s.ObjectiveValues[1]
		assert.LessOrEqual(t, sum, 10.0+1e-5)
		assert.GreaterOrEqual(t, sum, 0.0-1e-5)
	}
}

func Test_ParetoWeightedSum3Objectives(t *testing.T) {
	x := lp.NewVariable("x", lp.LpCategoryContinuous)
	y := lp.NewVariable("y", lp.LpCategoryContinuous)
	z := lp.NewVariable("z", lp.LpCategoryContinuous)
	variables := []lp.LpVariable{x, y, z}

	// Objectives 1..3: maximize x,y,z
	obj1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)})
	obj2 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, y)})
	obj3 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, z)})

	con := lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(1, x),
		lp.NewTerm(1, y),
		lp.NewTerm(1, z),
	})

	prog := lp.NewLinearProgram("Pareto Weighted Sum 3 Obj Test", variables)
	prog.AddObjective(lp.LpMaximise, obj1)
	prog.AddObjective(lp.LpMaximise, obj2)
	prog.AddObjective(lp.LpMaximise, obj3)
	prog.AddConstraint(con, lp.LpConstraintLE, 10)

	sols, err := SolvePareto(&prog, ParetoWeightedSum)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(sols), 3)

	for _, s := range sols {
		sum := s.ObjectiveValues[0] + s.ObjectiveValues[1] + s.ObjectiveValues[2]
		assert.IsClose(t, sum, 10.0, 1e-5)
	}
}

func Test_ParetoEpsilon3Objectives(t *testing.T) {
	x := lp.NewVariable("x", lp.LpCategoryContinuous)
	y := lp.NewVariable("y", lp.LpCategoryContinuous)
	z := lp.NewVariable("z", lp.LpCategoryContinuous)
	variables := []lp.LpVariable{x, y, z}

	obj1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)})
	obj2 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, y)})
	obj3 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, z)})

	con := lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(1, x),
		lp.NewTerm(1, y),
		lp.NewTerm(1, z),
	})

	prog := lp.NewLinearProgram("Pareto Epsilon 3 Obj Test", variables)
	prog.AddObjective(lp.LpMaximise, obj1)
	prog.AddObjective(lp.LpMaximise, obj2)
	prog.AddObjective(lp.LpMaximise, obj3)
	prog.AddConstraint(con, lp.LpConstraintLE, 10)

	sols, err := SolvePareto(&prog, ParetoEpsilonConstraint)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(sols), 3)

	for _, s := range sols {
		sum := s.ObjectiveValues[0] + s.ObjectiveValues[1] + s.ObjectiveValues[2]
		assert.LessOrEqual(t, sum, 10.0+1e-5)
		assert.GreaterOrEqual(t, sum, 0.0-1e-5)
	}
}

func Test_ParetoSingleObjectiveFails(t *testing.T) {
	x := lp.NewVariable("x", lp.LpCategoryContinuous)
	variables := []lp.LpVariable{x}

	obj := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)})
	con := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)})

	prog := lp.NewLinearProgram("Pareto Single Obj Test", variables)
	prog.AddObjective(lp.LpMaximise, obj)
	prog.AddConstraint(con, lp.LpConstraintLE, 10)

	_, err := SolvePareto(&prog, ParetoWeightedSum)
	assert.NotNil(t, err)
	_, err = SolvePareto(&prog, ParetoEpsilonConstraint)
	assert.NotNil(t, err)
}
