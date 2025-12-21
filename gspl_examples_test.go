package gspl_test

import (
	"fmt"

	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

func Example_readme() {
	variables := []lp.LpVariable{
		lp.NewVariable("x1"),
		lp.NewVariable("x2"),
		lp.NewVariable("x3"),
	}

	x1 := &variables[0]
	x2 := &variables[1]
	x3 := &variables[2]

	objective := lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(-6, *x1),
		lp.NewTerm(7, *x2),
		lp.NewTerm(4, *x3),
	})

	example := lp.NewLinearProgram("README Example", variables)
	example.AddObjective(lp.LpMinimise, objective)

	example.AddConstraint(lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(2, *x1),
		lp.NewTerm(5, *x2),
		lp.NewTerm(-1, *x3),
	}), lp.LpConstraintLE, 18)

	example.AddConstraint(lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(1, *x1),
		lp.NewTerm(-1, *x2),
		lp.NewTerm(-2, *x3),
	}), lp.LpConstraintLE, -14)

	example.AddConstraint(lp.NewExpression([]lp.LpTerm{
		lp.NewTerm(3, *x1),
		lp.NewTerm(2, *x2),
		lp.NewTerm(2, *x3),
	}), lp.LpConstraintEQ, 26)

	fmt.Printf("%s\n", example.String())

	sol, err := solver.Solve(&example)
	if err != nil {
		fmt.Println("solve error:", err)
		return
	}
	fmt.Printf("Optimal Objective Value: %.2f\n", sol.ObjectiveValue)
	fmt.Printf("Primal Solution: %v\n", sol.PrimalSolution.RawVector().Data)

	// Output:
	// README Example
	// Minimize: -6.00 * x1 + 7.00 * x2 + 4.00 * x3
	// Subject to:
	//   C1: 2.00 * x1 + 5.00 * x2 - 1.00 * x3 <= 18.000
	//   C2: -1.00 * x1 + 1.00 * x2 + 2.00 * x3 >= 14.000
	//   C3: 3.00 * x1 + 2.00 * x2 + 2.00 * x3 == 26.000
	//
	// Optimal Objective Value: 16.00
	// Primal Solution: [3 0 8.5]
}

func Example_minimisation() {
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

	fmt.Printf("%s\n", minProg.String())
	sol, err := solver.Solve(&minProg)
	if err != nil {
		fmt.Println("solve error:", err)
		return
	}
	fmt.Printf("Optimal Objective Value: %.2f\n", sol.ObjectiveValue)
	fmt.Printf("Primal Solution: %v\n", sol.PrimalSolution.RawVector().Data)
	// Output:
	// Minimisation Example
	// Minimize: 1.00 * x1 + 2.00 * x2 + 3.00 * x3 + 1.00 * x4 + 4.00 * x5
	// Subject to:
	//   C1: 1.00 * x1 + 1.00 * x2 + 1.00 * x3 + 1.00 * x4 + 1.00 * x5 >= 10.000
	//   C2: 1.00 * x1 + 2.00 * x2 + 1.00 * x3 <= 8.000
	//   C3: 1.00 * x2 + 1.00 * x4 + 1.00 * x5 <= 7.000
	//   C4: 1.00 * x1 + 1.00 * x3 + 1.00 * x5 >= 4.000
	//   C5: 1.00 * x4 + 1.00 * x5 <= 6.000
	//
	// Optimal Objective Value: 10.00
	// Primal Solution: [4 0 0 6 0]
}

func Example_iPMinimisation() {
	variables := []lp.LpVariable{
		lp.NewVariable("x1", lp.LpCategoryInteger),
		lp.NewVariable("x2", lp.LpCategoryInteger),
	}

	objTerms := []lp.LpTerm{
		lp.NewTerm(3, variables[0]),
		lp.NewTerm(2, variables[1]),
	}
	objective := lp.NewExpression(objTerms)

	con1Terms := []lp.LpTerm{
		lp.NewTerm(1.5, variables[0]),
		lp.NewTerm(1, variables[1]),
	}

	con2Terms := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(0.5, variables[1]),
	}

	lpProg := lp.NewLinearProgram("Non-Integer LP Example", variables)
	lpProg.AddObjective(lp.LpMinimise, objective)
	lpProg.AddConstraint(lp.NewExpression(con1Terms), lp.LpConstraintGE, 7)
	lpProg.AddConstraint(lp.NewExpression(con2Terms), lp.LpConstraintGE, 3)

	fmt.Printf("%s\n", lpProg.String())

	sol, err := solver.Solve(&lpProg)
	if err != nil {
		fmt.Println("solve error:", err)
		return
	}
	fmt.Printf("Optimal Objective Value: %.2f\n", sol.ObjectiveValue)
	fmt.Printf("Primal Solution: %v\n", sol.PrimalSolution.RawVector().Data)
	// Output:
	// Non-Integer LP Example
	// Minimize: 3.00 * x1 + 2.00 * x2
	// Subject to:
	//   C1: 1.50 * x1 + 1.00 * x2 >= 7.000
	//   C2: 1.00 * x1 + 0.50 * x2 >= 3.000
	// Integer variables: x1, x2
	//
	// Optimal Objective Value: 14.00
	// Primal Solution: [4 1]
}

func Example_maximisation() {
	variables := []lp.LpVariable{
		lp.NewVariable("x1", lp.LpCategoryContinuous),
		lp.NewVariable("x2", lp.LpCategoryContinuous),
	}

	objTerms := []lp.LpTerm{
		lp.NewTerm(5, variables[0]),
		lp.NewTerm(4, variables[1]),
	}
	objective := lp.NewExpression(objTerms)

	con1Terms := []lp.LpTerm{
		lp.NewTerm(2, variables[0]),
		lp.NewTerm(3, variables[1]),
	}
	con2Terms := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(1, variables[1]),
	}

	lpProg := lp.NewLinearProgram("Maximization Example", variables)
	lpProg.AddObjective(lp.LpMaximise, objective)
	lpProg.AddConstraint(lp.NewExpression(con1Terms), lp.LpConstraintLE, 12)
	lpProg.AddConstraint(lp.NewExpression(con2Terms), lp.LpConstraintLE, 5)

	fmt.Printf("%s\n", lpProg.String())

	sol, err := solver.Solve(&lpProg)
	if err != nil {
		fmt.Println("solve error:", err)
		return
	}
	fmt.Printf("Optimal Objective Value: %.2f\n", sol.ObjectiveValue)
	fmt.Printf("Primal Solution: %v\n", sol.PrimalSolution.RawVector().Data)
	// Output:
	// Maximization Example
	// Maximize: 5.00 * x1 + 4.00 * x2
	// Subject to:
	//   C1: 2.00 * x1 + 3.00 * x2 <= 12.000
	//   C2: 1.00 * x1 + 1.00 * x2 <= 5.000
	//
	// Optimal Objective Value: 25.00
	// Primal Solution: [5 0]
}

func Example_knapsackProblem() {
	variables := []lp.LpVariable{
		lp.NewVariable("x1", lp.LpCategoryInteger),
		lp.NewVariable("x2", lp.LpCategoryInteger),
		lp.NewVariable("x3", lp.LpCategoryInteger),
		lp.NewVariable("x4", lp.LpCategoryInteger),
		lp.NewVariable("x5", lp.LpCategoryInteger),
	}

	objTerms := []lp.LpTerm{
		lp.NewTerm(5, variables[0]),
		lp.NewTerm(3, variables[1]),
		lp.NewTerm(6, variables[2]),
		lp.NewTerm(6, variables[3]),
		lp.NewTerm(2, variables[4]),
	}
	objective := lp.NewExpression(objTerms)

	knapsackTerms := []lp.LpTerm{
		lp.NewTerm(1, variables[0]),
		lp.NewTerm(4, variables[1]),
		lp.NewTerm(7, variables[2]),
		lp.NewTerm(6, variables[3]),
		lp.NewTerm(2, variables[4]),
	}
	knapsackConstraint := lp.NewExpression(knapsackTerms)

	lpProg := lp.NewLinearProgram("Knapsack Problem", variables)
	lpProg.AddObjective(lp.LpMaximise, objective)
	lpProg.AddConstraint(knapsackConstraint, lp.LpConstraintLE, 15)

	for _, v := range variables {
		lpProg.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, v)}), lp.LpConstraintGE, 0)
		lpProg.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, v)}), lp.LpConstraintLE, 1)
	}

	fmt.Printf("%s\n", lpProg.String())

	sol, err := solver.Solve(&lpProg)
	if err != nil {
		fmt.Println("solve error:", err)
		return
	}
	fmt.Printf("Optimal Objective Value: %.2f\n", sol.ObjectiveValue)
	fmt.Printf("Primal Solution: %v\n", sol.PrimalSolution.RawVector().Data)
	// Output:
	// Knapsack Problem
	// Maximize: 5.00 * x1 + 3.00 * x2 + 6.00 * x3 + 6.00 * x4 + 2.00 * x5
	// Subject to:
	//   C1: 1.00 * x1 + 4.00 * x2 + 7.00 * x3 + 6.00 * x4 + 2.00 * x5 <= 15.000
	//   C2: 1.00 * x1 >= 0.000
	//   C3: 1.00 * x1 <= 1.000
	//   C4: 1.00 * x2 >= 0.000
	//   C5: 1.00 * x2 <= 1.000
	//   C6: 1.00 * x3 >= 0.000
	//   C7: 1.00 * x3 <= 1.000
	//   C8: 1.00 * x4 >= 0.000
	//   C9: 1.00 * x4 <= 1.000
	//   C10: 1.00 * x5 >= 0.000
	//   C11: 1.00 * x5 <= 1.000
	// Integer variables: x1, x2, x3, x4, x5
	//
	// Optimal Objective Value: 17.00
	// Primal Solution: [1 0 1 1 0]
}

func Example_transhipmentProblem() {
	cities := []int{1, 2, 3, 4}

	// --- Create edge variables x[i][j] for i != j ---
	x := make(map[[2]int]lp.LpVariable)
	for _, i := range cities {
		for _, j := range cities {
			if i != j {
				name := fmt.Sprintf("x%d%d", i, j)
				x[[2]int{i, j}] = lp.NewVariable(name, lp.LpCategoryInteger)
			}
		}
	}

	// --- Create MTZ variables u[2..4] ---
	u := make(map[int]lp.LpVariable)
	for i := 2; i <= 4; i++ {
		u[i] = lp.NewVariable(fmt.Sprintf("u%d", i), lp.LpCategoryContinuous)
	}

	var variables []lp.LpVariable
	for _, v := range x {
		variables = append(variables, v)
	}
	for _, v := range u {
		variables = append(variables, v)
	}

	prog := lp.NewLinearProgram("4-City TSP MTZ", variables)

	// --- Distance matrix ---
	dist := map[[2]int]float64{
		{1, 2}: 10, {1, 3}: 15, {1, 4}: 20,
		{2, 1}: 10, {2, 3}: 35, {2, 4}: 25,
		{3, 1}: 15, {3, 2}: 35, {3, 4}: 30,
		{4, 1}: 20, {4, 2}: 25, {4, 3}: 30,
	}

	// --- Objective ---
	var objTerms []lp.LpTerm
	for k, v := range x {
		objTerms = append(objTerms, lp.NewTerm(dist[k], v))
	}
	prog.AddObjective(lp.LpMinimise, lp.NewExpression(objTerms))

	// --- Outgoing edges constraints ---
	for _, i := range cities {
		var terms []lp.LpTerm
		for _, j := range cities {
			if i != j {
				terms = append(terms, lp.NewTerm(1, x[[2]int{i, j}]))
			}
		}
		prog.AddConstraint(lp.NewExpression(terms), lp.LpConstraintEQ, 1)
	}

	// --- Incoming edges constraints ---
	for _, j := range cities {
		var terms []lp.LpTerm
		for _, i := range cities {
			if i != j {
				terms = append(terms, lp.NewTerm(1, x[[2]int{i, j}]))
			}
		}
		prog.AddConstraint(lp.NewExpression(terms), lp.LpConstraintEQ, 1)
	}

	// --- MTZ subtour elimination constraints ---
	for i := range 2 {
		for j := range 2 {
			if i != j {
				terms := []lp.LpTerm{
					lp.NewTerm(1, u[i+2]),
					lp.NewTerm(-1, u[j+2]),
					lp.NewTerm(float64(len(cities)), x[[2]int{i + 2, j + 2}]),
				}
				prog.AddConstraint(lp.NewExpression(terms), lp.LpConstraintLE, float64(len(cities)-1))
			}
		}
	}

	// --- MTZ variable bounds ---
	for i := range 2 {
		prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, u[i+2])}), lp.LpConstraintGE, 2)
		prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, u[i+2])}), lp.LpConstraintLE, 4)
	}

	// --- Binary variable bounds ---
	for _, v := range x {
		prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, v)}), lp.LpConstraintGE, 0)
		prog.AddConstraint(lp.NewExpression([]lp.LpTerm{lp.NewTerm(1, v)}), lp.LpConstraintLE, 1)
	}

	// --- Solve ---
	sol, err := solver.Solve(&prog)
	if err != nil {
		fmt.Println("solve error:", err)
		return
	}

	fmt.Printf("Optimal Objective Value: %.2f\n", sol.ObjectiveValue)
	// Output:
	// Optimal Objective Value: 80.00
}
