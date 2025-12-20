package gmpl

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/chriso345/gspl/internal/lang/ast"
	"github.com/chriso345/gspl/lp"
)

func TestParseExampleFile(t *testing.T) {
	b, err := os.ReadFile("../testdata/example.gmpl")
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
	if !found["x1"] || !found["x2"] {
		t.Fatalf("expected x1 and x2 to be present, vars: %v", m.LP.Vars)
	}
	// objective should be maximise
	if m.LP.Sense != lp.LpMaximise {
		t.Fatalf("expected maximise sense, got %v", m.LP.Sense)
	}
}
