package solver

import (
	"context"
	"fmt"

	"github.com/chriso345/gspl/lp"
)

// Example_Solve demonstrates solving a simple LP using [Solve].
func ExampleSolve() {
	// Build a trivial LP: minimize x subject to x >= 1
	x := lp.NewVariable("x")
	p := lp.NewLinearProgram("Example LP", []lp.LpVariable{x})
	p.AddObjective(lp.LpMinimise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}), lp.LpConstraintGE, 1)

	// Run solve with a background context
	sol, err := Solve(&p, WithContext(context.Background()))
	if err != nil {
		fmt.Println("solve failed:", err)
		return
	}
	fmt.Printf("obj=%.0f, x=%.0f, status=%v\n", sol.ObjectiveValue, sol.PrimalSolution.AtVec(0), sol.Status)
	// Output: obj=1, x=1, status=Optimal
}

// ExampleNewSolverConfig demonstrates passing multiple [With*] options to [Solve].
func ExampleNewSolverConfig() {
	x := lp.NewVariable("x")
	p := lp.NewLinearProgram("t", []lp.LpVariable{x})
	p.AddObjective(lp.LpMinimise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}), lp.LpConstraintGE, 0)

	sol, _ := Solve(&p, WithTolerance(0.0001), WithMaxIterations(100), WithGapSensitivity(0.001), WithLogging(true))
	fmt.Printf("%.4f %v\n", sol.ObjectiveValue, sol.Status)
	// Output: 0.0000 Optimal
}

// ExampleWithContext demonstrates attaching a context via [WithContext] and using it with [Solve].
func ExampleWithContext() {
	x := lp.NewVariable("x")
	p := lp.NewLinearProgram("t", []lp.LpVariable{x})
	p.AddObjective(lp.LpMinimise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}), lp.LpConstraintGE, 0)

	sol, _ := Solve(&p, WithContext(context.Background()))
	fmt.Println(sol != nil)
	// Output: true
}

// ExampleWithTolerance demonstrates using [WithTolerance] with [Solve].
func ExampleWithTolerance() {
	// trivial LP minimize x subject to x >= 0
	x := lp.NewVariable("x")
	p := lp.NewLinearProgram("t", []lp.LpVariable{x})
	p.AddObjective(lp.LpMinimise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}), lp.LpConstraintGE, 0)

	sol, _ := Solve(&p, WithTolerance(1e-5))
	fmt.Printf("%.5f\n", sol.ObjectiveValue)
	// Output: 0.00000
}

// ExampleWithMaxIterations demonstrates using [WithMaxIterations] with [Solve].
func ExampleWithMaxIterations() {
	x := lp.NewVariable("x")
	p := lp.NewLinearProgram("t", []lp.LpVariable{x})
	p.AddObjective(lp.LpMinimise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}), lp.LpConstraintGE, 0)

	sol, _ := Solve(&p, WithMaxIterations(10))
	fmt.Println(sol.Status)
	// Output: Optimal
}

// ExampleWithGapSensitivity demonstrates using [WithGapSensitivity] with [Solve].
func ExampleWithGapSensitivity() {
	x := lp.NewVariable("x")
	p := lp.NewLinearProgram("t", []lp.LpVariable{x})
	p.AddObjective(lp.LpMinimise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}), lp.LpConstraintGE, 0)

	sol, _ := Solve(&p, WithGapSensitivity(0.5))
	fmt.Printf("%.1f\n", sol.ObjectiveValue)
	// Output: 0.0
}

// ExampleWithLogging demonstrates using [WithLogging] with [Solve].
func ExampleWithLogging() {
	x := lp.NewVariable("x")
	p := lp.NewLinearProgram("t", []lp.LpVariable{x})
	p.AddObjective(lp.LpMinimise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
	p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}), lp.LpConstraintGE, 0)

	sol, _ := Solve(&p, WithLogging(true))
	fmt.Println(sol.Status)
	// Output: Optimal
}
