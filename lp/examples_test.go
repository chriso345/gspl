package lp

import (
	"fmt"
)

// Example_NewLinearProgram demonstrates creating a [LinearProgram] and adding an objective and constraint.
func ExampleNewLinearProgram() {
	x := NewVariable("x")
	y := NewVariable("y")
	p := NewLinearProgram("Example LP", []LpVariable{x, y})
	p.AddObjective(LpMinimise, NewExpression([]LpTerm{NewTerm(1, x), NewTerm(2, y)}))
	p.AddConstraint(NewExpression([]LpTerm{NewTerm(1, x), NewTerm(1, y)}), LpConstraintGE, 3)

	fmt.Println(p.Description)
	// Output: Example LP
}

// Example_NewVariable demonstrates creating variables with different [LpCategory] values.
func ExampleNewVariable() {
	x := NewVariable("x")
	y := NewVariable("y", LpCategoryInteger)
	z := NewVariable("z", LpCategoryBinary)
	fmt.Printf("%s %d %s %d %s %d\n", x.Name, int(x.Category), y.Name, int(y.Category), z.Name, int(z.Category))
	// Output: x 0 y 1 z 2
}

// ExampleNewTerm demonstrates [NewTerm].
func ExampleNewTerm() {
	x := NewVariable("x")
	t := NewTerm(3, x)
	fmt.Printf("%.0f %s\n", t.Coefficient, t.Variable.Name)
	// Output: 3 x
}

// ExampleNewExpression demonstrates [NewExpression].
func ExampleNewExpression() {
	x := NewVariable("x")
	t := NewTerm(2, x)
	expr := NewExpression([]LpTerm{t})
	fmt.Printf("%d\n", len(expr.Terms))
	// Output: 1
}

// ExampleLinearProgram_AddObjective shows [LinearProgram].AddObjective negating coefficients for [LpMaximise].
func ExampleLinearProgram_AddObjective() {
	x := NewVariable("x")
	y := NewVariable("y")
	p := NewLinearProgram("o", []LpVariable{x, y})
	p.AddObjective(LpMaximise, NewExpression([]LpTerm{NewTerm(2, x), NewTerm(-3, y)}))
	// Coefficients are negated for maximisation
	fmt.Printf("%.0f %.0f %t\n", p.Objective.AtVec(0), p.Objective.AtVec(1), p.ObjectiveIsNegated)
	// Output: -2 3 true
}

// ExampleLinearProgram_AddConstraint shows adding a GE constraint and the slack/surplus column.
func ExampleLinearProgram_AddConstraint() {
	x := NewVariable("x")
	y := NewVariable("y")
	p := NewLinearProgram("c", []LpVariable{x, y})
	p.AddObjective(LpMinimise, NewExpression([]LpTerm{NewTerm(1, x), NewTerm(1, y)}))
	p.AddConstraint(NewExpression([]LpTerm{NewTerm(1, x), NewTerm(2, y)}), LpConstraintGE, 3)
	// Constraint row: 1 2 -1 (surplus) and objective expanded to length 3
	fmt.Printf("%.0f %.0f %.0f %d\n", p.Constraints.At(0, 0), p.Constraints.At(0, 1), p.Constraints.At(0, 2), p.Objective.Len())
	// Output: 1 2 -1 3
}
