package ast

import "testing"

func TestPositionString(t *testing.T) {
	p := Position{Line: 1, Column: 2}
	if s := p.String(); s != "1:2" {
		t.Fatalf("unexpected position string: %s", s)
	}
	p.SourceURI = "file"
	if s := p.String(); s != "file:1:2" {
		t.Fatalf("unexpected position string with URI: %s", s)
	}
}

func TestModuleNodeType(t *testing.T) {
	m := &Module{Name: "m"}
	if nt := m.NodeType(); nt != "Module" {
		t.Fatalf("unexpected node type: %s", nt)
	}
}
