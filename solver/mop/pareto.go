package mop

import (
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
	"gonum.org/v1/gonum/mat"
)

// SolvePareto solves the given multi-objective linear program using a Pareto
// approach. This method seeks to find a set of solutions that represent the best
// trade-offs among the objectives, known as the Pareto front. It returns a slice
// of MopSolution, each representing a non-dominated solution in the objective space.
func SolvePareto(
	prog *lp.LinearProgram,
	method ParetoMethod,
	opts ...solver.SolverOption,
) ([]*MopSolution, error) {

	if prog.Objective == nil {
		return nil, errors.New(errors.ErrInvalidInput, "primary objective must be defined", nil)
	}

	p, objs, err := cloneWithObjectives(prog)
	if err != nil {
		return nil, err
	}

	if len(objs) < 2 {
		return nil, errors.New(errors.ErrInvalidInput, "pareto requires >=2 objectives", nil)
	}

	switch method {
	case ParetoWeightedSum:
		return solveParetoWeightedSum(&p, objs, opts...)
	case ParetoEpsilonConstraint:
		return solveParetoEpsilon(&p, objs, opts...)
	default:
		return nil, errors.New(errors.ErrInvalidInput, "unknown pareto method", nil)
	}
}

func solveParetoWeightedSum(
	p *lp.LinearProgram,
	objs []*mat.VecDense,
	sopts ...solver.SolverOption,
) ([]*MopSolution, error) {
	cfg := solver.NewSolverConfig(sopts...)

	n := len(objs)
	if n == 0 {
		return nil, nil
	}

	weights := cfg.WeightedSums
	if len(weights) == 0 {
		weights = defaultSimplexWeights(n)
	}

	results := make([]*MopSolution, 0, len(weights))

	for _, w := range weights {
		if len(w) != n {
			return nil, errors.New(errors.ErrInvalidInput, "weight vector dimension mismatch", nil)
		}

		comb := mat.NewVecDense(objs[0].Len(), nil)
		for j, obj := range objs {
			for i := 0; i < obj.Len(); i++ {
				comb.SetVec(i, comb.AtVec(i)+w[j]*obj.AtVec(i))
			}
		}

		p.Objective = comb
		sol, err := solver.Solve(p, with(*cfg))
		if err != nil {
			return nil, err
		}

		vals := evaluateObjectives(sol.PrimalSolution, objs, p.ObjectiveIsNegated)
		results = append(results, &MopSolution{
			ObjectiveValues: vals,
			PrimalSolution:  sol.PrimalSolution,
			Status:          sol.Status,
		})
	}

	return uniqueSolutions(results), nil
}

func solveParetoEpsilon(
	p *lp.LinearProgram,
	objs []*mat.VecDense,
	sopts ...solver.SolverOption,
) ([]*MopSolution, error) {
	cfg := solver.NewSolverConfig(sopts...)

	if len(objs) < 2 {
		return nil, errors.New(errors.ErrInvalidInput, "epsilon-constraint requires >=2 objectives", nil)
	}

	primary := objs[0]
	others := objs[1:]

	p.Objective = primary
	base, err := solver.Solve(p, with(*cfg))
	if err != nil {
		return nil, err
	}

	results := []*MopSolution{
		{
			ObjectiveValues: evaluateObjectives(base.PrimalSolution, objs, p.ObjectiveIsNegated),
			PrimalSolution:  base.PrimalSolution,
			Status:          base.Status,
		},
	}

	epsilonSteps := cfg.EpsilonSteps
	for _, sec := range others {
		min, max := objectiveBounds(p, sec, sopts...)
		step := (max - min) / float64(epsilonSteps)

		for i := 1; i <= epsilonSteps; i++ {
			eps := min + float64(i)*step
			addEpsilonConstraint(p, sec, eps)

			sol, err := solver.Solve(p, with(*cfg))
			if err != nil {
				continue
			}

			results = append(results, &MopSolution{
				ObjectiveValues: evaluateObjectives(sol.PrimalSolution, objs, p.ObjectiveIsNegated),
				PrimalSolution:  sol.PrimalSolution,
				Status:          sol.Status,
			})
		}
	}

	return uniqueSolutions(results), nil
}

func evaluateObjectives(x *mat.VecDense, objs []*mat.VecDense, neg bool) []float64 {
	vals := make([]float64, len(objs))
	for j, obj := range objs {
		sum := 0.0
		for i := 0; i < x.Len() && i < obj.Len(); i++ {
			c := obj.AtVec(i)
			if neg {
				c = -c
			}
			sum += c * x.AtVec(i)
		}
		vals[j] = sum
	}
	return vals
}
