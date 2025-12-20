package solver

import (
	"testing"

	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/lp"

	"github.com/chriso345/gore/assert"
	"gonum.org/v1/gonum/mat"
)

// Test solving a simple non-integer LP
func TestSolve_LinearNoIP(t *testing.T) {
	// Minimize:  x + y
	// Subject to: x + y >= 4

	objective := mat.NewVecDense(2, []float64{1.0, 1.0})
	constraints := mat.NewDense(1, 2, []float64{1, 1})
	rhs := mat.NewVecDense(1, []float64{4})

	prog := &lp.LinearProgram{
		Sense:       lp.LpMinimise,
		Objective:   objective,
		Constraints: constraints,
		RHS:         rhs,
		ConTypes:    []lp.LpConstraintType{lp.LpConstraintGE},
		Vars: []lp.LpVariable{
			{Name: "x"},
			{Name: "y"},
		},
		Status: common.SolverStatusNotSolved,
	}

	t.Logf("[LinearNoIP] Starting solve...")
	sol, err := Solve(prog)
	if sol != nil {
		prog.ObjectiveValue = sol.ObjectiveValue
		prog.PrimalSolution = sol.PrimalSolution
		prog.Status = sol.Status
	}
	t.Logf("[LinearNoIP] Objective Value: %f", prog.ObjectiveValue)
	t.Logf("[LinearNoIP] Primal Solution: %v", prog.PrimalSolution.RawVector().Data)
	t.Logf("[LinearNoIP] Solver Status: %v", prog.Status)

	assert.Nil(t, err)
	assert.NotZero(t, prog.ObjectiveValue)
	assert.Equal(t, prog.PrimalSolution.Len(), 2)
}

// Test maximization sign flip
func TestSolve_Maximization(t *testing.T) {
	// Maximize: 3x
	// Subject to: x <= 10

	objective := mat.NewVecDense(1, []float64{3.0})
	constraints := mat.NewDense(1, 1, []float64{1})
	rhs := mat.NewVecDense(1, []float64{10})

	prog := &lp.LinearProgram{
		Sense:       lp.LpMaximise,
		Objective:   objective,
		Constraints: constraints,
		RHS:         rhs,
		ConTypes:    []lp.LpConstraintType{lp.LpConstraintLE},
		Vars:        []lp.LpVariable{{Name: "x"}},
		Status:      common.SolverStatusNotSolved,
	}

	t.Logf("[Maximization] Starting solve...")
	sol, err := Solve(prog)
	if sol != nil {
		prog.ObjectiveValue = sol.ObjectiveValue
		prog.PrimalSolution = sol.PrimalSolution
		prog.Status = sol.Status
	}
	t.Logf("[Maximization] Objective Value: %f", prog.ObjectiveValue)
	t.Logf("[Maximization] Primal Solution: %v", prog.PrimalSolution.RawVector().Data)
	t.Logf("[Maximization] Solver Status: %v", prog.Status)

	assert.Nil(t, err)
	assert.Equal(t, prog.ObjectiveValue, 30.0)
}

// Test minimisation
func TestSolve_Minimization(t *testing.T) {
	// Minimize: 5x
	// Subject to: x >= 2

	objective := mat.NewVecDense(1, []float64{5.0})
	constraints := mat.NewDense(1, 1, []float64{1})
	rhs := mat.NewVecDense(1, []float64{2})

	prog := &lp.LinearProgram{
		Sense:       lp.LpMinimise,
		Objective:   objective,
		Constraints: constraints,
		RHS:         rhs,
		ConTypes:    []lp.LpConstraintType{lp.LpConstraintGE},
		Vars:        []lp.LpVariable{{Name: "x"}},
		Status:      common.SolverStatusNotSolved,
	}

	t.Logf("[Minimization] Starting solve...")
	sol, err := Solve(prog)
	if sol != nil {
		prog.ObjectiveValue = sol.ObjectiveValue
		prog.PrimalSolution = sol.PrimalSolution
		prog.Status = sol.Status
	}
	t.Logf("[Minimization] Objective Value: %f", prog.ObjectiveValue)
	t.Logf("[Minimization] Primal Solution: %v", prog.PrimalSolution.RawVector().Data)
	t.Logf("[Minimization] Solver Status: %v", prog.Status)

	assert.Nil(t, err)
	assert.Greater(t, prog.ObjectiveValue, 0.0)
}

// Test that integer variable triggers integer solving
func TestSolve_MinIntegerProgram(t *testing.T) {
	// Minimize x
	// Constraint: x >= 3 with x integer

	objective := mat.NewVecDense(1, []float64{1.0})
	constraints := mat.NewDense(1, 1, []float64{1})
	rhs := mat.NewVecDense(1, []float64{3})

	prog := &lp.LinearProgram{
		Sense:       lp.LpMinimise,
		Objective:   objective,
		Constraints: constraints,
		RHS:         rhs,
		ConTypes:    []lp.LpConstraintType{lp.LpConstraintGE},
		Vars:        []lp.LpVariable{{Name: "x", Category: lp.LpCategoryInteger}},
		Status:      common.SolverStatusNotSolved,
	}

	t.Logf("[MinInteger] Starting solve...")
	sol, err := Solve(prog)
	if sol != nil {
		prog.ObjectiveValue = sol.ObjectiveValue
		prog.PrimalSolution = sol.PrimalSolution
		prog.Status = sol.Status
	}
	t.Logf("[MinInteger] Objective Value: %f", prog.ObjectiveValue)
	t.Logf("[MinInteger] Primal Solution: %v", prog.PrimalSolution.RawVector().Data)
	t.Logf("[MinInteger] Solver Status: %v", prog.Status)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, prog.ObjectiveValue, 3.0)
	assert.Equal(t, prog.PrimalSolution.Len(), 1)
}

func TestSolve_MaxIntegerProgram(t *testing.T) {
	// Maximize x
	// Constraint: x <= 7 with x integer

	objective := mat.NewVecDense(1, []float64{1.0})
	constraints := mat.NewDense(1, 1, []float64{1})
	rhs := mat.NewVecDense(1, []float64{7})

	prog := &lp.LinearProgram{
		Sense:       lp.LpMaximise,
		Objective:   objective,
		Constraints: constraints,
		RHS:         rhs,
		ConTypes:    []lp.LpConstraintType{lp.LpConstraintLE},
		Vars:        []lp.LpVariable{{Name: "x", Category: lp.LpCategoryInteger}},
		Status:      common.SolverStatusNotSolved,
	}

	t.Logf("[MaxInteger] Starting solve...")
	sol, err := Solve(prog)
	if sol != nil {
		prog.ObjectiveValue = sol.ObjectiveValue
		prog.PrimalSolution = sol.PrimalSolution
		prog.Status = sol.Status
	}
	t.Logf("[MaxInteger] Objective Value: %f", prog.ObjectiveValue)
	t.Logf("[MaxInteger] Primal Solution: %v", prog.PrimalSolution.RawVector().Data)
	t.Logf("[MaxInteger] Solver Status: %v", prog.Status)

	assert.Nil(t, err)
	assert.Equal(t, prog.ObjectiveValue, 7.0)
	assert.Equal(t, prog.PrimalSolution.Len(), 1)
}

// Test tolerance cleanup of tiny floating-point values
func TestSolve_ToleranceZeroing(t *testing.T) {
	objective := mat.NewVecDense(1, []float64{1})
	constraints := mat.NewDense(1, 1, []float64{1})
	rhs := mat.NewVecDense(1, []float64{0})

	initialSol := mat.NewVecDense(1, []float64{1e-10})

	prog := &lp.LinearProgram{
		Sense:          lp.LpMinimise,
		Objective:      objective,
		Constraints:    constraints,
		RHS:            rhs,
		ConTypes:       []lp.LpConstraintType{lp.LpConstraintLE},
		Vars:           []lp.LpVariable{{Name: "x"}},
		PrimalSolution: initialSol,
		Status:         common.SolverStatusNotSolved,
	}

	t.Logf("[Tolerance] Starting solve...")
	solution, err := Solve(prog, WithTolerance(1e-6))
	if solution != nil {
		prog.ObjectiveValue = solution.ObjectiveValue
		prog.PrimalSolution = solution.PrimalSolution
		prog.Status = solution.Status
	}
	t.Logf("[Tolerance] Objective Value: %f", prog.ObjectiveValue)
	t.Logf("[Tolerance] Primal Solution: %v", prog.PrimalSolution.RawVector().Data)
	t.Logf("[Tolerance] Solver Status: %v", prog.Status)

	assert.Nil(t, err)
	assert.Zero(t, prog.PrimalSolution.AtVec(0))
}
