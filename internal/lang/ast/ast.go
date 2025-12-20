package ast

import "github.com/chriso345/gspl/lp"

// Node is the minimal shared AST node interface.
type Node interface {
	NodeType() string
}

// Module is a simple top-level container used by minimal adapters.
// It may optionally contain a linear program when parsing succeeded.
type Module struct {
	Name string
	LP   *lp.LinearProgram
}

func (m *Module) NodeType() string { return "Module" }
