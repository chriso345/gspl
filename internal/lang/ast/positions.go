package ast

import "fmt"

// Position describes a source location in an input file.
type Position struct {
	SourceURI string
	Line      int
	Column    int
	Offset    int64
}

func (p Position) String() string {
	if p.SourceURI == "" {
		return fmt.Sprintf("%d:%d", p.Line, p.Column)
	}
	return fmt.Sprintf("%s:%d:%d", p.SourceURI, p.Line, p.Column)
}
