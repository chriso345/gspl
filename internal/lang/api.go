package lang

import (
	"context"
	"io"

	"github.com/chriso345/gspl/internal/lang/ast"
)

// Language abstracts parsing and per-language helpers.
type Language interface {
	// Name returns the canonical name used for registry lookups.
	Name() string

	// Parse parses the provided source and returns a node from the shared
	// AST package or a language-specific AST that implements ast.Node.
	Parse(ctx context.Context, src io.Reader, opts ...ParseOption) (ast.Node, error)
}

// ParseOption configures parsing behavior.
type ParseOption func(*ParseOptions)

// ParseOptions holds optional parser settings.
type ParseOptions struct {
	SourceURI string
	MaxBytes  int64
}
