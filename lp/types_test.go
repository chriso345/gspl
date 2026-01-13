package lp

import "testing"

func TestNewTermVariableExpression(t *testing.T) {
	x := NewVariable("x")
	if x.Name != "x" || x.IsSlack || x.IsArtificial || x.Category != LpCategoryContinuous {
		t.Fatalf("unexpected variable: %+v", x)
	}

	xInt := NewVariable("y", LpCategoryInteger)
	if xInt.Category != LpCategoryInteger {
		t.Fatalf("expected integer category")
	}

	term := NewTerm(5, x)
	if term.Coefficient != 5.0 || term.Variable.Name != "x" {
		t.Fatalf("unexpected term: %+v", term)
	}

	expr := NewExpression([]LpTerm{term})
	if len(expr.Terms) != 1 {
		t.Fatalf("expected 1 term")
	}
}

func TestNewVariablePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when passing multiple categories")
		}
	}()
	// Passing more than one category should panic
	NewVariable("x", LpCategoryInteger, LpCategoryBinary)
}

func TestLpStatusString(t *testing.T) {
	want := []string{"Not Solved", "Optimal", "Infeasible", "Unbounded"}
	for i, w := range want {
		if LpStatus(i).String() != w {
			t.Fatalf("unexpected string for status %d: got %q want %q", i, LpStatus(i).String(), w)
		}
	}
}
