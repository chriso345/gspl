package common

// SolverStatus represents the status of the solver
type SolverStatus int

const (
	SolverStatusNotSolved SolverStatus = iota
	SolverStatusOptimal
	SolverStatusInfeasible
	SolverStatusUnbounded
)

// String returns the string representation of the SolverStatus
func (s SolverStatus) String() string {
	switch s {
	case SolverStatusNotSolved:
		return "Not Solved"
	case SolverStatusOptimal:
		return "Optimal"
	case SolverStatusInfeasible:
		return "Infeasible"
	case SolverStatusUnbounded:
		return "Unbounded"
	default:
		return "Unknown"
	}
}

type VarCategory int

const (
	VarCategoryContinuous VarCategory = iota
	VarCategoryInteger
	VarCategoryBinary
)
