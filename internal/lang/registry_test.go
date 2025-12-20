package lang

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/chriso345/gspl/internal/lang/ast"
)

// testLang is a tiny Language implementation used by tests.
type testLang struct{}

func (t *testLang) Name() string { return "gmpl" }

func (t *testLang) Parse(ctx context.Context, src io.Reader, opts ...ParseOption) (ast.Node, error) {
	b, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	return &ast.Module{Name: string(b)}, nil
}

func TestParseBytesAndFile(t *testing.T) {
	ctx := context.Background()
	// register a lightweight test language to avoid import cycles
	RegisterLanguage(&testLang{})
	defer UnregisterLanguage("gmpl")

	b, err := os.ReadFile("testdata/example.gmpl")
	if err != nil {
		t.Fatalf("read example: %v", err)
	}
	node, err := ParseBytes(ctx, "gmpl", b)
	if err != nil {
		t.Fatalf("ParseBytes: %v", err)
	}
	m, ok := node.(*ast.Module)
	if !ok {
		t.Fatalf("unexpected node type: %T", node)
	}
	if !strings.Contains(m.Name, "maximize z") {
		t.Fatalf("unexpected content: %q", m.Name)
	}

	node2, err := ParseFile(ctx, "gmpl", "testdata/example.gmpl")
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	m2, ok := node2.(*ast.Module)
	if !ok {
		t.Fatalf("unexpected node type from ParseFile: %T", node2)
	}
	if m2.Name != m.Name {
		t.Fatalf("ParseFile != ParseBytes")
	}
}
