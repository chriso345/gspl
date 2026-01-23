package mps

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/chriso345/gspl/internal/lang/ast"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

func TestParseExampleFile(t *testing.T) {
	b, err := os.ReadFile("../testdata/example.mps")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}

	p := New()
	node, err := p.Parse(context.Background(), strings.NewReader(string(b)))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	m, ok := node.(*ast.Module)
	if !ok {
		t.Fatalf("unexpected node type: %T", node)
	}
	if m.LP == nil {
		t.Fatalf("expected LP to be non-nil")
	}
	// ensure primary variables exist
	found := map[string]bool{}
	for _, v := range m.LP.Vars {
		found[v.Name] = true
	}
	if !found["X1"] || !found["X2"] {
		t.Fatalf("expected X1 and X2 to be present, vars: %v", m.LP.Vars)
	}
	// objective should be minimise for MPS default
	if m.LP.Sense != lp.LpMinimise {
		t.Fatalf("expected minimise sense, got %v", m.LP.Sense)
	}

	// Solve and validate optimal solution
	sol, err := solver.Solve(m.LP)
	if err != nil {
		t.Fatalf("solve: %v", err)
	}
	if sol == nil {
		t.Fatalf("expected non-nil solution")
	}
	// For the example MPS: expect x1=0 (bounded by up 4, but objective min -> choose 0) and x2=10 (from LIM2 >=10)
	if sol.PrimalSolution.Len() < 2 {
		t.Fatalf("unexpected solution length: %d", sol.PrimalSolution.Len())
	}
	x1 := sol.PrimalSolution.AtVec(0)
	x2 := sol.PrimalSolution.AtVec(1)
	if diff := x1 - 0.0; diff < -1e-6 || diff > 1e-6 {
		t.Fatalf("unexpected X1: %v", x1)
	}
	if diff := x2 - 10.0; diff < -1e-6 || diff > 1e-6 {
		t.Fatalf("unexpected X2: %v", x2)
	}
	if diff := sol.ObjectiveValue - 20.0; diff < -1e-6 || diff > 1e-6 {
		t.Fatalf("unexpected objective: %v", sol.ObjectiveValue)
	}
}
