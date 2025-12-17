package solver_bench

import (
	"testing"

	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

func BenchmarkSolve_Small(b *testing.B) {
	for b.Loop() {
		variables := []lp.LpVariable{
			lp.NewVariable("x1", lp.LpCategoryInteger),
			lp.NewVariable("x2", lp.LpCategoryInteger),
		}

		objTerms := []lp.LpTerm{lp.NewTerm(3, variables[0]), lp.NewTerm(2, variables[1])}
		objective := lp.NewExpression(objTerms)

		con1Terms := []lp.LpTerm{lp.NewTerm(1.5, variables[0]), lp.NewTerm(1, variables[1])}
		con2Terms := []lp.LpTerm{lp.NewTerm(1, variables[0]), lp.NewTerm(0.5, variables[1])}

		prog := lp.NewLinearProgram("Small Benchmark", variables)
		prog.AddObjective(lp.LpMinimise, objective)
		prog.AddConstraint(lp.NewExpression(con1Terms), lp.LpConstraintGE, 7)
		prog.AddConstraint(lp.NewExpression(con2Terms), lp.LpConstraintGE, 3)

		if _, err := solver.Solve(&prog); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSolve_Medium(b *testing.B) {
	for b.Loop() {
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

		if _, err := solver.Solve(&prog); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSolve_Large(b *testing.B) {
	for b.Loop() {
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

		if _, err := solver.Solve(&prog); err != nil {
			b.Fatal(err)
		}
	}
}
