package lp

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestPrintSolution(t *testing.T) {
	x1 := NewVariable("x1")
	x2 := NewVariable("x2")
	lp := NewLinearProgram("Test LP", []LpVariable{x1, x2})

	lp.AddObjective(LpMinimise, NewExpression([]LpTerm{
		NewTerm(2, x1),
		NewTerm(3, x2),
	}))

	lp.AddConstraint(NewExpression([]LpTerm{
		NewTerm(1, x1),
		NewTerm(1, x2),
	}), LpConstraintLE, 5)

	lp.PrimalSolution = mat.NewVecDense(3, []float64{1, 2, 0})
	lp.ObjectiveValue = 8.0

	// Redirect stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lp.PrintSolution()

	if err := w.Close(); err != nil {
		t.Fatalf("failed to close pipe writer: %v", err)
	}
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	os.Stdout = oldStdout

	if buf.Len() == 0 {
		t.Fatalf("expected output from PrintSolution")
	}
	out := buf.String()
	if !strings.Contains(out, "Test LP") || !strings.Contains(out, "ObjectiveValue") || !strings.Contains(out, "x1") || !strings.Contains(out, "x2") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestLinearProgramStringUndefinedObjective(t *testing.T) {
	p := NewLinearProgram("desc", []LpVariable{NewVariable("x1")})
	s := p.String()
	if !strings.Contains(s, "<undefined objective>") {
		t.Fatalf("expected undefined objective in string, got: %s", s)
	}
}

func TestLinearProgramStringWithObjectiveAndConstraints(t *testing.T) {
	x := NewVariable("x1")
	p := NewLinearProgram("desc", []LpVariable{x})
	p.AddObjective(LpMinimise, NewExpression([]LpTerm{NewTerm(-2, x)}))
	p.AddConstraint(NewExpression([]LpTerm{NewTerm(1, x)}), LpConstraintLE, 3)
	p.Vars[1].Category = LpCategoryInteger // mark slack as integer to hit integer listing
	p.Vars[0].Category = LpCategoryBinary  // mark original var as binary
	s := p.String()
	if !strings.Contains(s, "Minimize") || !strings.Contains(s, "Subject to:") || !strings.Contains(s, "Integer variables") || !strings.Contains(s, "Binary variables") {
		t.Fatalf("unexpected string output: %s", s)
	}
}
