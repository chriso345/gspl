package lp

import (
	"strings"
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/internal/common"
)

func TestNewLinearProgram(t *testing.T) {
	vars := []LpVariable{
		NewVariable("x1"),
		NewVariable("x2", LpCategoryInteger),
	}
	lp := NewLinearProgram("Test LP", vars)

	assert.Equal(t, lp.Description, "Test LP")
	assert.Equal(t, len(lp.Vars), 2)
	assert.Equal(t, lp.Status, common.SolverStatusNotSolved)
}

func TestLinearProgramStringMaximiseAndSigns(t *testing.T) {
	x := NewVariable("x1")
	y := NewVariable("x2")
	p := NewLinearProgram("desc", []LpVariable{x, y})
	p.AddObjective(LpMaximise, NewExpression([]LpTerm{NewTerm(2, x), NewTerm(-3, y)}))
	p.AddConstraint(NewExpression([]LpTerm{NewTerm(1, x), NewTerm(-1, y)}), LpConstraintLE, 4)
	s := p.String()
	if !strings.Contains(s, "Maximize") {
		t.Fatalf("expected Maximize in string, got: %s", s)
	}
	// ensure signs and coefficients are present
	if !strings.Contains(s, "2.00 * x1") || !strings.Contains(s, "3.00 * x2") {
		t.Fatalf("unexpected objective formatting: %s", s)
	}
}
