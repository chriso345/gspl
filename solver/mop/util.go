package mop

import (
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
	"gonum.org/v1/gonum/mat"
)

func cloneWithObjectives(prog *lp.LinearProgram) (lp.LinearProgram, []*mat.VecDense, error) {
	// Detect inconsistent objective direction changed after AddObjective
	if prog.ObjectiveIsNegated && prog.Sense != lp.LpMaximise {
		return lp.LinearProgram{}, nil,
			errors.New(errors.ErrInvalidInput, "inconsistent objective direction", nil)
	}
	if !prog.ObjectiveIsNegated && prog.Sense != lp.LpMinimise {
		return lp.LinearProgram{}, nil,
			errors.New(errors.ErrInvalidInput, "inconsistent objective direction", nil)
	}

	// Shallow copy
	p := *prog

	// Deep copy mutable fields
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

	// Ensure constraints exist
	if p.Constraints == nil {
		p.Constraints = mat.NewDense(1, len(p.Vars), nil)
		p.RHS = mat.NewVecDense(1, []float64{0})
	}

	// Collect objective vectors in stored form
	objs := make([]*mat.VecDense, 0, 1+len(p.SecondaryObjectives))
	objs = append(objs, mat.VecDenseCopyOf(p.Objective))

	for _, sec := range p.SecondaryObjectives {
		if sec == nil {
			objs = append(objs, mat.NewVecDense(len(p.Vars), nil))
			continue
		}
		if sec.Len() < len(p.Vars) {
			ext := mat.NewVecDense(len(p.Vars), nil)
			for i := 0; i < sec.Len(); i++ {
				ext.SetVec(i, sec.AtVec(i))
			}
			sec = ext
		}
		objs = append(objs, sec)
	}

	return p, objs, nil
}

func objectiveBounds(
	p *lp.LinearProgram,
	obj *mat.VecDense,
	opts ...solver.SolverOption,
) (min, max float64) {

	origObj := p.Objective
	defer func() { p.Objective = origObj }()

	// Ensure objective vector matches number of variables
	if obj.Len() < len(p.Vars) {
		ext := mat.NewVecDense(len(p.Vars), nil)
		for i := 0; i < obj.Len(); i++ {
			ext.SetVec(i, obj.AtVec(i))
		}
		p.Objective = ext
	} else if obj.Len() > len(p.Vars) {
		// truncate if longer
		ext := mat.NewVecDense(len(p.Vars), nil)
		for i := 0; i < len(p.Vars); i++ {
			ext.SetVec(i, obj.AtVec(i))
		}
		p.Objective = ext
	} else {
		p.Objective = obj
	}

	// Minimise
	p.Sense = lp.LpMinimise
	minSol, err := solver.Solve(p, opts...)
	if err != nil {
		return 0, 0
	}
	min = minSol.ObjectiveValue

	// Maximise
	p.Sense = lp.LpMaximise
	maxSol, err := solver.Solve(p, opts...)
	if err != nil {
		return min, min
	}
	max = maxSol.ObjectiveValue

	return min, max
}

func addEpsilonConstraint(
	p *lp.LinearProgram,
	obj *mat.VecDense,
	epsilon float64,
) {
	terms := make([]lp.LpTerm, 0, obj.Len())

	for i := 0; i < obj.Len() && i < len(p.Vars); i++ {
		coef := obj.AtVec(i)
		if coef == 0 {
			continue
		}
		if p.ObjectiveIsNegated {
			coef = -coef
		}
		terms = append(terms, lp.NewTerm(coef, p.Vars[i]))
	}

	expr := lp.NewExpression(terms)

	dir := lp.LpConstraintLE
	if p.Sense == lp.LpMaximise {
		dir = lp.LpConstraintGE
	}

	p.AddConstraint(expr, dir, epsilon)
}

func uniqueSolutions(sols []*MopSolution) []*MopSolution {
	const tol = 1e-9

	uniq := make([]*MopSolution, 0, len(sols))
	for _, s := range sols {
		dup := false
		for _, u := range uniq {
			if almostEqual(s.ObjectiveValues, u.ObjectiveValues) {
				dup = true
				break
			}
		}
		if !dup {
			uniq = append(uniq, s)
		}
	}

	// Sort by first objective for stable ordering
	for i := 0; i < len(uniq); i++ {
		for j := i + 1; j < len(uniq); j++ {
			if uniq[j].ObjectiveValues[0] < uniq[i].ObjectiveValues[0] {
				uniq[i], uniq[j] = uniq[j], uniq[i]
			}
		}
	}

	return uniq
}

func almostEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	const tol = 1e-9
	for i := range a {
		d := a[i] - b[i]
		if d < 0 {
			d = -d
		}
		if d > tol {
			return false
		}
	}
	return true
}

func defaultSimplexWeights(n int) [][]float64 {
	weights := make([][]float64, n)
	for i := 0; i < n; i++ {
		w := make([]float64, n)
		w[i] = 1
		weights[i] = w
	}
	return weights
}
