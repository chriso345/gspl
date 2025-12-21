package lang

import (
	"context"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/chriso345/gspl/internal/lang/ast"
)

type testLang struct{}

func (t *testLang) Name() string { return "gmpl" }

func (t *testLang) Parse(ctx context.Context, src io.Reader, opts ...ParseOption) (ast.Node, error) {
	b, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	return &ast.Module{Name: string(b)}, nil
}

type fakeLang struct{ name string }

func (f *fakeLang) Name() string { return f.name }
func (f *fakeLang) Parse(ctx context.Context, src io.Reader, opts ...ParseOption) (ast.Node, error) {
	return &ast.Module{Name: f.name}, nil
}

func registerLang(t *testing.T, l Language) {
	t.Helper()
	if err := RegisterLanguage(l); err != nil {
		t.Fatalf("RegisterLanguage(%q): %v", l.Name(), err)
	}
	t.Cleanup(func() { UnregisterLanguage(l.Name()) })
}

func TestParseBytesAndFile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	registerLang(t, &testLang{})

	b, err := os.ReadFile("testdata/example.gmpl")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	t.Run("ParseBytes", func(t *testing.T) {
		node, err := ParseBytes(ctx, "gmpl", b)
		if err != nil {
			t.Fatalf("ParseBytes: %v", err)
		}
		m := mustModule(t, node)
		if !strings.Contains(m.Name, "maximize z") {
			t.Fatalf("unexpected content: %q", m.Name)
		}
	})

	t.Run("ParseFile", func(t *testing.T) {
		node, err := ParseFile(ctx, "gmpl", "testdata/example.gmpl")
		if err != nil {
			t.Fatalf("ParseFile: %v", err)
		}
		m := mustModule(t, node)
		if !strings.Contains(m.Name, "maximize z") {
			t.Fatalf("unexpected content: %q", m.Name)
		}
	})
}

func TestRegisterAndUnregisterLanguage(t *testing.T) {
	t.Parallel()

	l := &fakeLang{name: "fake"}
	UnregisterLanguage("fake")

	if err := RegisterLanguage(l); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Cleanup(func() { UnregisterLanguage("fake") })

	if err := RegisterLanguage(l); err != nil {
		t.Fatalf("idempotent register failed: %v", err)
	}

	if err := RegisterLanguage(&fakeLang{name: "fake"}); err == nil {
		t.Fatal("expected error on conflicting registration")
	}
}

func TestParseUnsupportedLanguage(t *testing.T) {
	t.Parallel()

	if node, err := Parse(context.Background(), "no-such", nil); err == nil {
		t.Fatalf("expected error, got node=%v", node)
	}
}

func TestParseFileErrors(t *testing.T) {
	t.Parallel()

	t.Run("file not found", func(t *testing.T) {
		if _, err := ParseFile(context.Background(), "no-such", "definitely_missing.txt"); err == nil {
			t.Fatal("expected error for missing file")
		}
	})

	t.Run("unsupported language via ParseBytes", func(t *testing.T) {
		if _, err := ParseBytes(context.Background(), "no-such", []byte("x")); err == nil {
			t.Fatal("expected error for unsupported language")
		}
	})
}

func TestOptionsAndPanics(t *testing.T) {
	t.Parallel()

	t.Run("options do not panic", func(t *testing.T) {
		opts := &ParseOptions{}
		WithSourceURI("/tmp/foo")(opts)
		WithMaxBytes(10)(opts)
	})

	t.Run("MustRegisterLanguage panics on conflict", func(t *testing.T) {
		l := &fakeLang{name: "panicfake"}
		registerLang(t, l)

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic from MustRegisterLanguage")
			}
		}()

		MustRegisterLanguage(&fakeLang{name: "panicfake"})
	})
}

func TestParseFileSuccessWithFakeLang(t *testing.T) {
	t.Parallel()

	f, err := os.CreateTemp("", "langtest")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	registerLang(t, &fakeLang{name: "filefake"})

	if _, err := ParseFile(context.Background(), "filefake", f.Name()); err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}
}

func mustModule(t *testing.T, node ast.Node) *ast.Module {
	t.Helper()
	m, ok := node.(*ast.Module)
	if !ok {
		t.Fatalf("unexpected node type: %T", node)
	}
	return m
}
