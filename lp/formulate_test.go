package lp

import (
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestAddObjective(t *testing.T) {
	x1 := NewVariable("x1")
	x2 := NewVariable("x2")
	lp := NewLinearProgram("Test LP", []LpVariable{x1, x2})

	expr := NewExpression([]LpTerm{
		NewTerm(3, x1),
		NewTerm(5, x2),
	})

	lp.AddObjective(LpMinimise, expr)
	assert.Equal(t, lp.Objective.AtVec(0), 3.0)
	assert.Equal(t, lp.Objective.AtVec(1), 5.0)

	lpMax := NewLinearProgram("Max LP", []LpVariable{x1, x2})
	lpMax.AddObjective(LpMaximise, expr)
	assert.Equal(t, lpMax.Objective.AtVec(0), -3.0)
	assert.Equal(t, lpMax.Objective.AtVec(1), -5.0)
}

func TestAddObjectivePanic(t *testing.T) {
	lp := NewLinearProgram("Test LP", []LpVariable{})
	expr := NewExpression([]LpTerm{
		NewTerm(1, NewVariable("x1")),
	})

	assert.Panic(t, func() {
		lp.AddObjective(LpMinimise, expr)
	})
}

func TestAddConstraintWithoutObjectivePanic(t *testing.T) {
	x1 := NewVariable("x1")
	lp := NewLinearProgram("Test LP", []LpVariable{x1})
	expr := NewExpression([]LpTerm{
		NewTerm(1, x1),
	})
	assert.Panic(t, func() {
		lp.AddConstraint(expr, LpConstraintLE, 5)
	})
}

func TestAddConstraint(t *testing.T) {
	x1 := NewVariable("x1")
	x2 := NewVariable("x2")
	lp := NewLinearProgram("Test LP", []LpVariable{x1, x2})
	lp.AddObjective(LpMinimise, NewExpression([]LpTerm{
		NewTerm(1, x1),
		NewTerm(1, x2),
	}))

	expr := NewExpression([]LpTerm{
		NewTerm(1, x1),
		NewTerm(2, x2),
	})
	lp.AddConstraint(expr, LpConstraintLE, 4)

	assert.Equal(t, lp.Constraints.RawMatrix().Rows, 1)
	assert.Equal(t, len(lp.Vars), 3) // slack variable added

	// GE constraint with negative RHS flips
	expr2 := NewExpression([]LpTerm{
		NewTerm(2, x1),
		NewTerm(1, x2),
	})
	lp.AddConstraint(expr2, LpConstraintGE, -3)
	assert.Equal(t, lp.Constraints.RawMatrix().Rows, 2)
}

func TestAddConstraintFlipAndSlack(t *testing.T) {
	vars := []LpVariable{NewVariable("x1")}
	p := NewLinearProgram("desc", vars)
	// define an objective so AddConstraint can run
	p.AddObjective(LpMinimise, NewExpression([]LpTerm{NewTerm(1, vars[0])}))
	// Add a constraint with negative RHS which should be flipped and produce a surplus variable
	p.AddConstraint(NewExpression([]LpTerm{NewTerm(2, vars[0])}), LpConstraintLE, -5)
	if p.RHS.AtVec(0) != 5 {
		t.Fatalf("expected RHS 5 after flip, got %v", p.RHS.AtVec(0))
	}
	if p.ConTypes[0] != LpConstraintGE {
		t.Fatalf("expected constraint type GE after flip, got %v", p.ConTypes[0])
	}
	// slack variable should have been appended
	slack := p.Vars[len(p.Vars)-1]
	if !slack.IsSlack {
		t.Fatalf("expected last var to be slack")
	}
	// surplus for GE should be -1 in the current row
	if p.Constraints.At(0, p.Constraints.RawMatrix().Cols-1) != -1 {
		t.Fatalf("expected surplus -1, got %v", p.Constraints.At(0, p.Constraints.RawMatrix().Cols-1))
	}
}

func TestAddObjectiveVariableNotFoundPanic(t *testing.T) {
	p := NewLinearProgram("desc", []LpVariable{NewVariable("x1")})
	// expression refers to a variable not in p.Vars
	expr := NewExpression([]LpTerm{NewTerm(1, NewVariable("missing"))})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when objective references unknown variable")
		}
	}()
	p.AddObjective(LpMinimise, expr)
}
