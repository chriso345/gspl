package solver

import (
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/lp"
	"gonum.org/v1/gonum/mat"
)

// MopSolution contains the result of a multi-objective optimization.
type MopSolution struct {
	ObjectiveValues []float64
	PrimalSolution  *mat.VecDense
	Status          common.SolverStatus
}

// SolveLexicographic solves the given multi-objective linear program using
// a lexicographic approach. This method optimises the objectives in a predefined
// order, ensuring that the optimal value of each objective is achieved before
// moving on to the next. It returns a single MopSolution representing the optimal
// solution that respects the priority of objectives.
func SolveLexicographic(prog *lp.LinearProgram, opts ...SolverOption) (*MopSolution, error) {
	if prog.Objective == nil {
		return nil, errors.New(errors.ErrInvalidInput, "primary objective must be defined", nil)
	}

	// Detect inconsistent objective direction changed after AddObjective
	if prog.ObjectiveIsNegated && prog.Sense != lp.LpMaximise {
		return nil, errors.New(errors.ErrInvalidInput, "inconsistent objective direction", nil)
	}
	if !prog.ObjectiveIsNegated && prog.Sense != lp.LpMinimise {
		return nil, errors.New(errors.ErrInvalidInput, "inconsistent objective direction", nil)
	}

	p := *prog
	if prog.Objective != nil {
		p.Objective = mat.VecDenseCopyOf(prog.Objective)
	}
	if prog.Constraints != nil {
		p.Constraints = mat.DenseCopyOf(prog.Constraints)
	}
	if prog.RHS != nil {
		p.RHS = mat.VecDenseCopyOf(prog.RHS)
	}
	p.ConTypes = append([]lp.LpConstraintType(nil), prog.ConTypes...)
	p.Vars = append([]lp.LpVariable(nil), prog.Vars...)
	if prog.SecondaryObjectives != nil {
		p.SecondaryObjectives = make([]*mat.VecDense, len(prog.SecondaryObjectives))
		for i, v := range prog.SecondaryObjectives {
			if v != nil {
				p.SecondaryObjectives[i] = mat.VecDenseCopyOf(v)
			}
		}
	}

	if p.Constraints == nil {
		p.Constraints = mat.NewDense(1, len(p.Vars), nil)
		p.RHS = mat.NewVecDense(1, []float64{0})
	}

	sol, err := Solve(&p, opts...)
	if err != nil {
		return nil, err
	}

	objVals := []float64{sol.ObjectiveValue}
	lastSol := sol
	if len(p.SecondaryObjectives) == 0 {
		return &MopSolution{ObjectiveValues: objVals, PrimalSolution: lastSol.PrimalSolution, Status: lastSol.Status}, nil
	}

	vecToExpr := func(vec *mat.VecDense) lp.LpExpression {
		terms := make([]lp.LpTerm, 0, vec.Len())
		for i := 0; i < vec.Len() && i < len(p.Vars); i++ {
			coef := vec.AtVec(i)
			if coef == 0 {
				continue
			}
			// If the stored objective has been negated to represent maximisation,
			// flip the coefficient back when forming constraints.
			if p.ObjectiveIsNegated {
				coef = -coef
			}
			terms = append(terms, lp.NewTerm(coef, p.Vars[i]))
		}
		return lp.NewExpression(terms)
	}

	dir := lp.LpConstraintLE
	if p.Sense == lp.LpMaximise {
		dir = lp.LpConstraintGE
	}

	p.AddConstraint(vecToExpr(p.Objective), dir, sol.ObjectiveValue)

	for _, sec := range p.SecondaryObjectives {
		// Ensure secondary vector matches number of variables
		if sec.Len() < len(p.Vars) {
			ext := mat.NewVecDense(len(p.Vars), nil)
			for i := 0; i < sec.Len(); i++ {
				ext.SetVec(i, sec.AtVec(i))
			}
			sec = ext
		}
		// Ensure secondary is stored in the same negation state as the primary
		if p.ObjectiveIsNegated {
			for i := 0; i < sec.Len(); i++ {
				sec.SetVec(i, -sec.AtVec(i))
			}
		}
		p.Objective = sec
		lastSol, err = Solve(&p, opts...)
		if err != nil {
			return nil, err
		}
		objVals = append(objVals, lastSol.ObjectiveValue)
		p.AddConstraint(vecToExpr(p.Objective), dir, lastSol.ObjectiveValue)
	}

	return &MopSolution{ObjectiveValues: objVals, PrimalSolution: lastSol.PrimalSolution, Status: lastSol.Status}, nil
}

// SolvePareto solves the given multi-objective linear program using a Pareto
// approach. This method seeks to find a set of solutions that represent the best
// trade-offs among the objectives, known as the Pareto front. It returns a slice
// of MopSolution, each representing a non-dominated solution in the objective space.
func SolvePareto(prog *lp.LinearProgram, opts ...SolverOption) ([]*MopSolution, error) {
	return nil, nil
}
