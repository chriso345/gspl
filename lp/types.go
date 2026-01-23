package lp

import "github.com/chriso345/gspl/internal/common"

// LpExpression represents the left-hand side of a linear expression and
// holds a slice of [LpTerm].
type LpExpression struct {
	Terms []LpTerm
}

// NewExpression creates a new [LpExpression] with the given terms.
func NewExpression(terms []LpTerm) LpExpression {
	return LpExpression{terms}
}

// LpTerm represents a term in a linear expression, consisting of a coefficient and a [LpVariable].
type LpTerm struct {
	Coefficient float64
	Variable    LpVariable // These get added to the variable list in the LinearProgram??
}

// NewTerm returns a new [LpTerm] with the given coefficient and [LpVariable].
func NewTerm(coefficient float64, variable LpVariable) LpTerm {
	return LpTerm{coefficient, variable}
}

// LpVariable represents a variable in a linear programming problem and
// includes metadata such as slack/artificial flags and its [LpCategory].
type LpVariable struct {
	Name         string
	IsSlack      bool
	IsArtificial bool
	Category     LpCategory
}

// NewVariable creates a new [LpVariable] with the given name and optional [LpCategory].
func NewVariable(name string, category ...LpCategory) LpVariable {
	if len(category) > 1 {
		panic("Only one LpCategory can be specified for a variable")
	}
	if len(category) == 0 {
		return LpVariable{name, false, false, LpCategoryContinuous} // Default to continuous variable
	}
	return LpVariable{name, false, false, category[0]}
}

// LpCategory is an alias for [common.VarCategory] and classifies variables
// (continuous, integer, binary).
type LpCategory = common.VarCategory

const (
	LpCategoryContinuous LpCategory = iota
	LpCategoryInteger
	LpCategoryBinary
)

// LpSense represents the problem sense for a [LinearProgram].
// Possible values are [LpMinimise] and [LpMaximise].
type LpSense int

const (
	LpMinimise LpSense = iota
	LpMaximise
)

// LpStatus represents the current status of solving a [LinearProgram].
// Values include [LpStatusNotSolved], [LpStatusOptimal], [LpStatusInfeasible], and [LpStatusUnbounded].
type LpStatus int

const (
	LpStatusNotSolved LpStatus = iota
	LpStatusOptimal
	LpStatusInfeasible
	LpStatusUnbounded
)

// String returns the string representation of the LpStatus.
func (s LpStatus) String() string {
	return [...]string{
		"Not Solved",
		"Optimal",
		"Infeasible",
		"Unbounded",
	}[s]
}

// LpConstraintType represents the relation used in a constraint (<=, =, >=).
// Constants are [LpConstraintLE], [LpConstraintEQ], and [LpConstraintGE].
type LpConstraintType int

const (
	LpConstraintLE LpConstraintType = iota - 1
	LpConstraintEQ
	LpConstraintGE
)
