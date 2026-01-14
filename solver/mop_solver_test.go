package solver

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/lp"
	"gonum.org/v1/gonum/mat"
)

func TestSolveLexicographic_ExistingObjectives(t *testing.T) {
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
	objective2 := lp.NewExpression(terms2)

	prog := lp.NewLinearProgram("Lexicographic Minimisation Example", variables)

	_, err := SolveLexicographic(&prog)
	assert.NotNil(t, err)

	prog.AddObjective(lp.LpMinimise, objective)

	_, err = SolveLexicographic(&prog)
	assert.Nil(t, err)

	// Add secondary objective explicitly
	// helper to convert expression to vector
	exprToVec := func(expr lp.LpExpression) *mat.VecDense {
		v := mat.NewVecDense(len(prog.Vars), nil)
		for _, t := range expr.Terms {
			for i, vv := range prog.Vars {
				if vv.Name == t.Variable.Name {
					v.SetVec(i, t.Coefficient)
					break
				}
			}
		}
		return v
	}
	prog.SecondaryObjectives = append(prog.SecondaryObjectives, mat.VecDenseCopyOf(prog.Objective))
	prog.SecondaryObjectives = append(prog.SecondaryObjectives, exprToVec(objective2))
	res, err := SolveLexicographic(&prog)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(res.ObjectiveValues))
}

func TestSolveLexicographic_DoesNotMutateOriginalProgram(t *testing.T) {
	vars := []lp.LpVariable{lp.NewVariable("a"), lp.NewVariable("b")}
	obj1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0]), lp.NewTerm(0, vars[1])})
	obj2 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0]), lp.NewTerm(1, vars[1])})
	prog := lp.NewLinearProgram("No Mutation", vars)
	prog.AddObjective(lp.LpMinimise, obj1)
	prog.AddObjective(lp.LpMinimise, obj2)

	origObj := mat.VecDenseCopyOf(prog.Objective)

	_, err := SolveLexicographic(&prog)
	assert.Nil(t, err)
	// original program objective should be unchanged
	for i := 0; i < origObj.Len(); i++ {
		assert.Equal(t, origObj.AtVec(i), prog.Objective.AtVec(i))
	}
}

func TestSolveLexicographic_MinFullSolve(t *testing.T) {
	// Minimise two objectives lexicographically
	// vars x,y >=0
	vars := []lp.LpVariable{lp.NewVariable("x"), lp.NewVariable("y")}
	// primary: minimise x + 10y
	obj1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0]), lp.NewTerm(10, vars[1])})
	// secondary: minimise x
	_ = lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0]), lp.NewTerm(0, vars[1])})

	prog := lp.NewLinearProgram("Min Lex", vars)
	// constraints: x + y >= 10  and x <= 8
	prog.AddObjective(lp.LpMinimise, obj1)
	prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0]), lp.NewTerm(1, vars[1])}), lp.LpConstraintGE, 10)
	prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0])}), lp.LpConstraintLE, 8)
	// add secondary
	v1 := mat.NewVecDense(2, nil)
	v1.SetVec(0, 1)
	prog.SecondaryObjectives = append(prog.SecondaryObjectives, v1)

	res, err := SolveLexicographic(&prog)
	assert.Nil(t, err)
	// primary optimal: minimise x+10y subject to x+y>=10 and x<=8 -> best is x=8,y=2 -> obj1=8+20=28
	assert.Equal(t, 28.0, res.ObjectiveValues[0])
	// secondary should then minimise x given primary fixed -> x=8
	assert.Equal(t, 8.0, res.ObjectiveValues[1])
}

func TestSolveLexicographic_MaxFullSolve(t *testing.T) {
	// Maximise two objectives lexicographically
	vars := []lp.LpVariable{lp.NewVariable("x"), lp.NewVariable("y")}
	// primary: maximise x + y
	obj1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0]), lp.NewTerm(1, vars[1])})
	// secondary: maximise y
	_ = lp.NewExpression([]lp.LpTerm{lp.NewTerm(0, vars[0]), lp.NewTerm(1, vars[1])})

	prog := lp.NewLinearProgram("Max Lex", vars)
	prog.AddObjective(lp.LpMaximise, obj1)
	// constraints: x + 2y <= 10, x <= 6
	prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0]), lp.NewTerm(2, vars[1])}), lp.LpConstraintLE, 10)
	prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0])}), lp.LpConstraintLE, 6)
	// add secondary
	v2 := mat.NewVecDense(2, nil)
	v2.SetVec(1, 1)
	prog.SecondaryObjectives = append(prog.SecondaryObjectives, v2)

	res, err := SolveLexicographic(&prog)
	assert.Nil(t, err)
	// Ensure solution is optimal and secondary objective equals expected
	assert.Equal(t, common.SolverStatusOptimal, res.Status)
	assert.Equal(t, 2.0, res.ObjectiveValues[1])
	if res.PrimalSolution != nil {
		assert.True(t, res.PrimalSolution.AtVec(0) >= 0)
		assert.True(t, res.PrimalSolution.AtVec(1) >= 0)
	}
}

func TestSolveLexicographic_MixedDirectionsFails(t *testing.T) {
	vars := []lp.LpVariable{lp.NewVariable("x"), lp.NewVariable("y")}
	obj1 := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, vars[0])})
	prog := lp.NewLinearProgram("Mixed", vars)
	prog.AddObjective(lp.LpMinimise, obj1)
	v3 := mat.NewVecDense(2, nil)
	v3.SetVec(1, 1)
	prog.SecondaryObjectives = append(prog.SecondaryObjectives, v3)
	// Now flip program sense to simulate inconsistent directions
	prog.Sense = lp.LpMaximise
	_, err := SolveLexicographic(&prog)
	assert.NotNil(t, err)
}
