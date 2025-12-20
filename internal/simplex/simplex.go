package simplex

import (
	"math"

	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
	"github.com/chriso345/gspl/internal/matrix"
	"gonum.org/v1/gonum/mat"
)

func Simplex(scf *common.StandardComputationalForm, config *common.SolverConfig) error {
	m, n := scf.Constraints.Dims()
	sm := &simplexMethod{
		m: m,
		n: n,
	}

	I := matrix.Eye(m)

	// Phase 1: Set up the auxilary problem
	sm.A = mat.NewDense(m, n+m, nil)
	for i := range m {
		for j := range n {
			sm.A.Set(i, j, scf.Constraints.At(i, j))
		}
		for j := range m {
			sm.A.Set(i, n+j, I.At(i, j))
		}
	}

	// Construct the cost vector for Phase 1, [0,...,0,1,...,1]
	sm.c = mat.NewVecDense(n+m, nil)
	for i := range n + m {
		if i < n {
			sm.c.SetVec(i, 0.)
		} else {
			sm.c.SetVec(i, 1.)
		}
	}

	sm.cb = mat.NewVecDense(m, nil)
	for i := range m {
		sm.cb.SetVec(i, float64(n+i)) // Artificial variables as initial basis
	}

	sm.B = matrix.ExtractColumns(sm.A, sm.cb)
	sm.b = scf.RHS

	// Keep original constraints pointer so we can detect changes later (cheap check)
	origConstraints := scf.Constraints

	// Run Phase 1 of the RSM
	err := RSM(sm, 1, config)
	if err != nil {
		return errors.New(errors.ErrNumericalFailure, "error in Phase 1 of Simplex", err)
	}

	// Check infeasibility
	if sm.flag == common.SolverStatusOptimal && sm.value > config.Tolerance {
		*scf.Status = common.SolverStatusInfeasible
		return nil
	}

	// Phase 2: Set up the original problem
	// Remove artificial variables from basis if their xb = 0
	err = removeArtificialFromBasis(sm)
	if err != nil {
		*scf.Status = common.SolverStatusInfeasible
		return errors.New(errors.ErrNumericalFailure, "error removing artificial variables from basis", err)
	}

	// This _should_ always be true, but just in case
	if scf.Constraints != origConstraints {
		sm.A = mat.NewDense(m, n+m, nil)
		for i := range m {
			for j := range n {
				sm.A.Set(i, j, scf.Constraints.At(i, j))
			}
			for j := range m {
				sm.A.Set(i, n+j, I.At(i, j))
			}
		}
	}

	sm.c = mat.NewVecDense(n+m, nil)
	for i := range n {
		sm.c.SetVec(i, scf.Objective.AtVec(i))
	}
	for i := range m {
		sm.c.SetVec(n+i, 0.)
	}

	sm.cb = sm.indices
	sm.B = matrix.ExtractColumns(sm.A, sm.cb)
	// Reuse RHS from Phase 1 (sm.b already points to scf.RHS) to avoid extra allocations
	sm.b = scf.RHS

	// Run Phase 2 of the RSM
	err = RSM(sm, 2, config)
	if err != nil {
		return errors.New(errors.ErrNumericalFailure, "error in Phase 2 of Simplex", err)
	}
	*scf.Status = sm.flag
	if sm.flag == common.SolverStatusOptimal {
		*scf.ObjectiveValue = sm.value
		scf.PrimalSolution = sm.x
		// debug removed
	}

	return nil
}

func RSM(sm *simplexMethod, phase int, config *common.SolverConfig) error {
	_maxIter := 1000 // Simple safeguard

	n := sm.n
	if phase == 1 {
		n += sm.m
	}

	// Initialise the rsmResult (populate embedded fields)
	sm.flag = common.SolverStatusNotSolved
	sm.value = 0. // z
	sm.x = mat.NewVecDense(n, nil)
	sm.indices = sm.cb

	// Initialise other variables
	B := sm.B // RSM mutates B in-place and owns the basis; updateB writes into this matrix
	cb := mat.NewVecDense(sm.m, nil)
	for i := range sm.m {
		index := int(sm.indices.AtVec(i))
		cb.SetVec(i, sm.c.AtVec(index))
	}

	for range _maxIter {
		xb := mat.NewVecDense(sm.m, nil) // Basic solution
		err := xb.SolveVec(B, sm.b)
		if err != nil {
			// Basis is singular, return error
			return errors.New(errors.ErrNumericalFailure, "error solving for basic solution", err)
		}

		// Finding the leaving variable
		var BT mat.Dense
		BT.CloneFrom(B.T())

		sm.pi = mat.NewVecDense(sm.m, nil) // Dual variables
		err = sm.pi.SolveVec(&BT, cb)
		if err != nil {
			// Basis is singular, return error
			return errors.New(errors.ErrNumericalFailure, "error solving for dual variables", err)
		}

		fe := enteringVariable{
			A:       sm.A,
			pi:      sm.pi,
			c:       sm.c,
			isbasic: mat.NewVecDense(n, nil),

			epsilon: config.Tolerance,
		}

		for i := range sm.m {
			index := int(sm.indices.AtVec(i))
			if index < n {
				fe.isbasic.SetVec(index, 1.)
			}
		}

		err = findEnter(&fe)
		if err != nil {
			return errors.New(errors.ErrNumericalFailure, "error finding entering variable", err)
		}

		if fe.s == -1 {
			// Optimal solution found
			sm.flag = common.SolverStatusOptimal
			sm.value = 0.
			for i := range sm.m {
				index := int(sm.indices.AtVec(i))
				sm.x.SetVec(index, xb.AtVec(i))
				sm.value += cb.AtVec(i) * xb.AtVec(i)
			}
			return nil
		}

		// Finding the leaving variable
		fl := leavingVariable{
			B:       B,
			indices: sm.indices,
			as:      fe.as,
			xb:      xb,
			phase:   phase,
			n:       n,
		}

		err = findLeave(&fl)
		if err != nil {
			return errors.New(errors.ErrNumericalFailure, "error finding leaving variable", err)
		}

		if fl.r == -1 {
			// Unbounded solution
			sm.flag = common.SolverStatusUnbounded
			// Set primal solution vector to zero
			sm.x = mat.NewVecDense(n, nil)
			sm.value = 0.
			return nil
		}

		// Update B, cb, and indices
		bu := basisUpdate{
			BMat:    B,
			indices: sm.indices,
			cb:      cb,
			as:      fe.as,
			s:       fe.s,
			r:       fl.r,
			cs:      fe.cs,
		}

		if err := updateB(&bu); err != nil {
			return errors.New(errors.ErrNumericalFailure, "error updating basis", err)
		}

	}
	return errors.New(errors.ErrNumericalFailure, "max iterations reached in RSM", nil)
}

func findEnter(fe *enteringVariable) error {
	fe.s = -1
	fe.cs = 0.
	minrc := math.Inf(1)
	tol := -fe.epsilon

	n := fe.isbasic.Len()
	m, _ := fe.A.Dims()

	// Reuse or allocate the 'as' vector once per call
	if fe.as == nil {
		fe.as = mat.NewVecDense(m, nil)
	}

	for j := range n {
		if fe.isbasic.AtVec(j) == 0 {
			// Compute dot product without allocating a temporary vector
			dot := 0.0
			for i := range m {
				dot += fe.pi.AtVec(i) * fe.A.At(i, j)
			}
			rc := fe.c.AtVec(j) - dot

			if rc < minrc {
				minrc = rc
				fe.s = j
				fe.cs = fe.c.AtVec(j)

				// Reuse the preallocated as vector
				for i := range m {
					fe.as.SetVec(i, fe.A.At(i, j))
				}
			}
		}
	}

	if minrc >= tol {
		fe.s = -1
		for i := range m {
			fe.as.SetVec(i, 0)
		}
		fe.cs = 0.
	}

	return nil
}

func findLeave(fl *leavingVariable) error {
	fl.r = -1

	var Binv mat.Dense
	if err := Binv.Inverse(fl.B); err != nil {
		return errors.New(errors.ErrNumericalFailure, "error inverting basis matrix", err)
	}

	directionVec := mat.NewVecDense(fl.as.Len(), nil)
	directionVec.MulVec(&Binv, fl.as)

	m := fl.xb.Len()
	theta := math.Inf(1)

	for i := range m {
		dirVal := directionVec.AtVec(i)
		indexVal := int(fl.indices.AtVec(i))

		if fl.phase == 2 && indexVal > fl.n {
			if dirVal != 0 {
				fl.r = i
				return nil
			}
		} else {
			if dirVal > 0 {
				ratio := fl.xb.AtVec(i) / dirVal
				if ratio < theta {
					theta = ratio
					fl.r = i
				}
			}
		}
	}

	return nil
}

func updateB(bu *basisUpdate) error {
	m, _ := bu.BMat.Dims()

	for i := range m {
		bu.BMat.Set(i, bu.r, bu.as.AtVec(i))
	}

	bu.indices.SetVec(bu.r, float64(bu.s))
	bu.cb.SetVec(bu.r, bu.cs)

	return nil
}

// removeArtificialFromBasis removes artificial variables from the basis before Phase 2.
// If an artificial variable has a positive value, it returns an error (infeasible LP).
func removeArtificialFromBasis(sm *simplexMethod) error {
	for i := 0; i < sm.m; i++ {
		index := int(sm.indices.AtVec(i))
		if index >= sm.n { // artificial variable
			if math.Abs(sm.x.AtVec(index)) < 1e-8 {
				// Replace with a non-basic original variable
				replaced := false
				for j := 0; j < sm.n; j++ {
					if !contains(sm.indices, j) {
						sm.indices.SetVec(i, float64(j))
						replaced = true
						break
					}
				}
				if !replaced {
					return errors.New(errors.ErrInfeasible, "cannot remove artificial variable from basis: no non-basic original variable available", nil)
				}
			} else {
				return errors.New(errors.ErrInfeasible, "LP is infeasible: artificial variable in basis with positive value", nil)
			}
		}
	}
	return nil
}

func contains(vec *mat.VecDense, val int) bool {
	for i := 0; i < vec.Len(); i++ {
		if int(vec.AtVec(i)) == val {
			return true
		}
	}
	return false
}
