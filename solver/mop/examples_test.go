package mop

import (
	"fmt"

	"github.com/chriso345/gspl/lp"
	"gonum.org/v1/gonum/mat"
)

func ExampleSolveLexicographic() {
	x := lp.NewVariable("x")
	y := lp.NewVariable("y")
	p := lp.NewLinearProgram("Lex", []lp.LpVariable{x, y})
	// primary: minimise x + 10y
	p.AddObjective(lp.LpMinimise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x), lp.NewTerm(10, y)}))
	// constraints: x + y >= 10  and x <= 8
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x), lp.NewTerm(1, y)}), lp.LpConstraintGE, 10)
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}), lp.LpConstraintLE, 8)
	// secondary: minimise x
	v := lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)})
	p.SecondaryObjectives = append(p.SecondaryObjectives, matCopyVec(v, p.Vars))

	sol, err := SolveLexicographic(&p)
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Printf("primary=%.0f secondary=%.0f status=%v\n", sol.ObjectiveValues[0], sol.ObjectiveValues[1], sol.Status)
	// Output: primary=28 secondary=8 status=Optimal
}

func ExampleSolvePareto_weightedSum() {
	x := lp.NewVariable("x")
	y := lp.NewVariable("y")
	p := lp.NewLinearProgram("Pareto", []lp.LpVariable{x, y})
	p.AddObjective(lp.LpMaximise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddObjective(lp.LpMaximise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, y)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x), lp.NewTerm(1, y)}), lp.LpConstraintLE, 10)

	// Provide simple weights to sample the front
	sols, err := SolvePareto(&p, ParetoWeightedSum, WithWeightedSums([][]float64{{1, 0}, {0, 1}}))
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println(len(sols) >= 2)
	// Output: true
}

func ExampleSolvePareto_epsilonConstraint() {
	x := lp.NewVariable("x")
	y := lp.NewVariable("y")
	p := lp.NewLinearProgram("Pareto", []lp.LpVariable{x, y})
	p.AddObjective(lp.LpMaximise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddObjective(lp.LpMaximise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, y)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x), lp.NewTerm(1, y)}), lp.LpConstraintLE, 10)

	// Provide epsilon steps to sample the front
	sols, err := SolvePareto(&p, ParetoEpsilonConstraint, WithEpsilonSteps(3))
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Println(len(sols) >= 2)
	// Output: true
}

// helper to copy an expression into a *mat.VecDense of length n
func matCopyVec(expr lp.LpExpression, vars []lp.LpVariable) *mat.VecDense {
	v := mat.NewVecDense(len(vars), nil)
	for _, t := range expr.Terms {
		for i, vv := range vars {
			if vv.Name == t.Variable.Name {
				v.SetVec(i, t.Coefficient)
				break
			}
		}
	}
	return v
}
