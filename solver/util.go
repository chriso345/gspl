package solver

import (
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/lp"
)

// ErrorKind and Error are re-exported for public API use
type ErrorKind = errors.ErrorKind

var (
	ErrUnknown          = errors.ErrUnknown
	ErrInfeasible       = errors.ErrInfeasible
	ErrUnbounded        = errors.ErrUnbounded
	ErrNumericalFailure = errors.ErrNumericalFailure
	ErrInvalidInput     = errors.ErrInvalidInput
)

type Error = errors.Error

// hasIPConstraints checks if the linear program has any integer or binary constraints.
func hasIPConstraints(prog *lp.LinearProgram) bool {
	for _, v := range prog.Vars {
		if v.Category != lp.LpCategoryContinuous {
			return true
		}
	}
	return false
}
