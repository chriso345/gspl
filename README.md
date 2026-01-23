# gspl

**gspl** (short for **Go Simplex Programming Library**, pronounced *gospel*) is a lightweight solver for linear programming problems, built in Go using the revised simplex method. It emphasizes clarity, performance, and seamless integration into Go applications.

---

## Features

* Solves optimal, unbounded, infeasible, and degenerate linear problems.
* Implements an efficient revised simplex algorithm.
* Clean and idiomatic API for modeling and solving LPs.
* Focused on numerical stability and usability.
* Performs basic branch-and-bound techniques for pure integer problems.
* Provides native multi-objective optimization support.
* Builtin bindings for C and Python (coming soon).

`gspl` is ideal for embedding linear optimization into Go-based software.

---

## Installation

Ensure you have Go installed, then run:

```bash
go get github.com/chriso345/gspl
```

---

## Usage

Full examples are available in the [gspl_examples_test.go](gspl_examples_test.go) file and the [examples/](examples/) directory.

### Variables

Create decision variables as a slice:

```go
variables := []lp.LpVariable{
  lp.NewVariable("x1"),
  lp.NewVariable("x2"),
}
```

You can access and pass them by value or take their address when needed:

```go
x1 := &variables[0]
```

Forcing integer constraints is done when creating the variable:

```go
variables := []lp.LpVariable{
  lp.NewVariable("x1", lp.LpCategoryInteger),
  lp.NewVariable("x2", lp.LpCategoryInteger),
}
```

### Objective Function

Build the objective using terms:

```go
objective := lp.NewExpression([]lp.LpTerm{
  lp.NewTerm(5, variables[0]),
  lp.NewTerm(3, variables[1]),
})
```

Add it to the LP:

```go
prog := lp.NewLinearProgram("My LP", variables)
prog.AddObjective(lp.LpMaximise, objective)
```

### Constraints

Each constraint uses an expression, a comparison type, and a right-hand side:

```go
prog.AddConstraint(lp.NewExpression([]lp.LpTerm{
  lp.NewTerm(2, variables[0]),
  lp.NewTerm(3, variables[1]),
}), lp.LpConstraintLE, 10)
```

Constraint types:

* `gspl.LpConstraintLE` - less than or equal
* `gspl.LpConstraintGE` - greater than or equal
* `gspl.LpConstraintEQ` - equality

### Solving

```go
solution, err := solver.Solve(&prog)
if err != nil {
  // handle error
}
fmt.Printf("obj=%.2f, x=%.2f\n", solution.ObjectiveValue, solution.PrimalSolution.AtVec(0))
```

This solves the model and prints variable values and the objective result.

### Multi-Objective Support

gspl provides native multi-objective optimization capabilities via the `mop` package.

```go
x := lp.NewVariable("x")
y := lp.NewVariable("y")
p := lp.NewLinearProgram("Pareto", []lp.LpVariable{x, y})
p.AddObjective(lp.LpMaximise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x)}))
p.AddObjective(lp.LpMaximise, lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, y)}))
p.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, x), lp.NewTerm(1, y)}), lp.LpConstraintLE, 10)

sols, err := mop.SolvePareto(&p, mop.ParetoWeightedSum, mop.WithWeightedSums([][]float64{{1, 0}, {0, 1}}))
if err != nil {
  fmt.Println("err", err)
  return
}
```

---

## License

Licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.
