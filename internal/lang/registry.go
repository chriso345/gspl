package lang

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/chriso345/gspl/internal/lang/ast"
)

var (
	mu       sync.RWMutex
	registry = map[string]Language{}
)

// RegisterLanguage registers a language implementation. It is idempotent when
// the same instance is registered again.
func RegisterLanguage(l Language) error {
	mu.Lock()
	defer mu.Unlock()
	name := l.Name()
	if existing, ok := registry[name]; ok {
		if existing == l {
			return nil
		}
		return fmt.Errorf("language %q already registered", name)
	}
	registry[name] = l
	return nil
}

// MustRegisterLanguage registers a language and panics on error. Useful for init.
func MustRegisterLanguage(l Language) {
	if err := RegisterLanguage(l); err != nil {
		panic(err)
	}
}

// Parse looks up a registered language and delegates parsing to it.
func Parse(ctx context.Context, langName string, src io.Reader, opts ...ParseOption) (ast.Node, error) {
	mu.RLock()
	l, ok := registry[langName]
	mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", langName)
	}
	return l.Parse(ctx, src, opts...)
}

// ParseBytes is a convenience wrapper for parsing in-memory bytes.
func ParseBytes(ctx context.Context, langName string, b []byte, opts ...ParseOption) (ast.Node, error) {
	return Parse(ctx, langName, bytes.NewReader(b), opts...)
}

// ParseFile opens the named file and parses it.
func ParseFile(ctx context.Context, langName, filename string, opts ...ParseOption) (ast.Node, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(ctx, langName, f, opts...)
}

// UnregisterLanguage removes a language from the registry. Intended for tests.
func UnregisterLanguage(name string) {
	mu.Lock()
	defer mu.Unlock()
	delete(registry, name)
}
