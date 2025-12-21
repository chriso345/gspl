package solver

import (
	"context"
	"math"

	"github.com/chriso345/gspl/internal/brancher"
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/internal/simplex"
	"github.com/chriso345/gspl/lp"
	"gonum.org/v1/gonum/mat"
)

// Solution contains the results returned by Solve.
//
// The returned Solution provides a snapshot of the computed objective value,
// the primal solution vector, and a status code describing optimality,
// infeasibility, or unboundedness. The Solution contains copies of data
// that callers can safely inspect without referencing the original
// lp.LinearProgram.
//
// Note: for performance the solver may temporarily reuse internal SCF fields
// that point into the provided LinearProgram; therefore callers MUST NOT mutate
// the provided *lp.LinearProgram concurrently with a call to Solve.
// To cancel a long-running solve pass a context using the WithContext option.
type Solution struct {
	ObjectiveValue float64
	PrimalSolution *mat.VecDense
	Status         common.SolverStatus
}

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

// Solve solves the given linear program and returns a Solution and an error.
//
// The function returns a populated *Solution on success, or a non-nil error if
// the solve failed. Solve respects context cancellation when a context is
// provided via SolverOption (WithContext). It may temporarily link into fields
// of the provided LinearProgram for efficiency; therefore the provided program
// must not be mutated concurrently. Solve is safe to call concurrently as long
// as each goroutine uses a distinct *lp.LinearProgram.
func Solve(prog *lp.LinearProgram, opts ...SolverOption) (*Solution, error) {
	// Apply options
	options := NewSolverConfig(opts...)
	if options.Ctx == nil {
		options.Ctx = context.Background()
	}

	tol := options.Tolerance

	if hasIPConstraints(prog) {
		ip := newIP(prog)

		// Respect context cancellation
		select {
		case <-options.Ctx.Done():
			return nil, options.Ctx.Err()
		default:
		}

		// Call the Integer Programming solver
		err := brancher.BranchAndBound(ip, options)
		if err != nil {
			return nil, errors.New(errors.ErrUnknown, "integer solve failed", err)
		}

		sol := &Solution{Status: *ip.SCF.Status}
		sol.ObjectiveValue = ip.BestObj
		sol.PrimalSolution = mat.NewVecDense(ip.SCF.NumPrimals, nil)
		if ip.BestSolution != nil {
			// For integer programs, round the primal solution to integer values
			for i := 0; i < ip.SCF.NumPrimals; i++ {
				item := ip.BestSolution.AtVec(i)
				if item < tol && item > -tol {
					continue
				}
				rounded := math.Round(item)
				sol.PrimalSolution.SetVec(i, rounded)
			}
			// ip.BestObj is already stored in the original problem sense by the
			// branch-and-bound routine; use it rather than recomputing from the
			// possibly-negated lp.Objective vector.
		}

		return sol, nil
	}

	// Create the SCF instance
	scf := newSCF(prog)

	// Respect context cancellation
	select {
	case <-options.Ctx.Done():
		return nil, options.Ctx.Err()
	default:
	}

	// Call the Simplex solver
	if err := simplex.Simplex(scf, options); err != nil {
		return nil, errors.New(errors.ErrUnknown, "simplex failed", err)
	}

	// Copy the solution back without mutating the original problem state
	sol := &Solution{Status: *scf.Status}

	// Flip objective back to original sense for maximisation problems
	if scf.IsMaximization {
		sol.ObjectiveValue = -(*scf.ObjectiveValue)
	} else {
		sol.ObjectiveValue = *scf.ObjectiveValue
	}

	// Ensure we never dereference a nil PrimalSolution from the SCF
	sol.PrimalSolution = mat.NewVecDense(scf.NumPrimals, nil)
	if scf.PrimalSolution != nil {
		for i := 0; i < scf.NumPrimals; i++ {
			item := scf.PrimalSolution.AtVec(i)
			if item < tol && item > -tol {
				continue
			}
			sol.PrimalSolution.SetVec(i, item)
		}
	}

	return sol, nil
}

// newSCF creates a new SCF instance for the linear program
func newSCF(prog *lp.LinearProgram) *common.StandardComputationalForm {
	slackIndices := make([]int, len(prog.Vars))
	numPrimals := 0
	for i, constr := range prog.Vars {
		if constr.IsSlack {
			slackIndices[i] = i
		} else {
			slackIndices[i] = -1
			numPrimals++
		}
	}

	// Copy objective. If the program is a maximisation and the objective has not
	// already been negated by AddObjective, negate it here to convert to the
	// solver's minimisation form.
	objCopy := mat.VecDenseCopyOf(prog.Objective)
	if prog.Sense == lp.LpMaximise && !prog.ObjectiveIsNegated {
		objCopy.ScaleVec(-1, objCopy)
	}
	return &common.StandardComputationalForm{
		Objective:   objCopy,
		Constraints: prog.Constraints,
		RHS:         prog.RHS,

		PrimalSolution: prog.PrimalSolution,

		// Link back to the original problem
		ObjectiveValue: &prog.ObjectiveValue,
		Status:         &prog.Status,
		SlackIndices:   slackIndices,
		NumPrimals:     numPrimals,
		// Record original sense so results can be flipped back if needed
		IsMaximization: prog.Sense == lp.LpMaximise,
	}
}

// newIP creates a new IP instance for the linear program
func newIP(prog *lp.LinearProgram) *common.IntegerProgram {
	ip := &common.IntegerProgram{
		SCF: newSCF(prog),
	}
	// Initialize BestObj appropriately for minimisation/maximisation
	if prog.Sense == lp.LpMaximise {
		ip.BestObj = math.Inf(-1)
	} else {
		ip.BestObj = math.Inf(1)
	}
	return ip
}
