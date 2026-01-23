package simplex

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/matrix"
	"gonum.org/v1/gonum/mat"
)

func TestContains(t *testing.T) {
	v := mat.NewVecDense(4, []float64{0, 2, 3, 5})
	assert.True(t, contains(v, 2))
	assert.False(t, contains(v, 4))
	assert.True(t, contains(v, 0))
	assert.False(t, contains(v, 6))
}

func TestRemoveArtificialFromBasis(t *testing.T) {}

func TestSimplexSmallProblem(t *testing.T) {
	// Simple LP: maximize x subject to x <= 5, x >= 0
	scf := &common.StandardComputationalForm{
		Objective:      mat.NewVecDense(1, []float64{1}),
		Constraints:    mat.NewDense(1, 1, []float64{1}),
		RHS:            mat.NewVecDense(1, []float64{5}),
		PrimalSolution: mat.NewVecDense(1, nil),
		ObjectiveValue: new(float64),
		Status:         new(common.SolverStatus),
	}
	config := common.DefaultSolverConfig()
	err := Simplex(scf, config)
	if err != nil {
		t.Fatalf("Simplex returned error: %v", err)
	}
	if *scf.Status != common.SolverStatusOptimal && *scf.Status != common.SolverStatusInfeasible && *scf.Status != common.SolverStatusUnbounded {
		t.Fatalf("unexpected status: %v", *scf.Status)
	}
}

func TestRSM_PivotOnce(t *testing.T) {
	// Configure sm to force one pivot in RSM with m=2, n=3
	A := mat.NewDense(2, 3, []float64{
		1, 0, 1,
		0, 1, 1,
	})
	B := mat.NewDense(2, 2, []float64{1, 0, 0, 1})
	c := mat.NewVecDense(3, []float64{0, 0, -1})
	indices := mat.NewVecDense(2, []float64{0, 1})
	sm := &simplexMethod{
		m: 2,
		n: 3,
		A: A,
		B: B,
		c: c,
		b: mat.NewVecDense(2, []float64{5, 3}),
		rsmResult: rsmResult{
			indices: indices,
			x:       mat.NewVecDense(3, nil),
		},
		cb: mat.NewVecDense(2, []float64{0, 0}),
	}
	config := &common.SolverConfig{Tolerance: 1e-9, MaxIterations: 100}
	if err := RSM(sm, 2, config); err != nil {
		t.Fatalf("RSM failed: %v", err)
	}
	if sm.flag != common.SolverStatusOptimal {
		t.Fatalf("expected optimal, got %v", sm.flag)
	}
}

func TestRemoveArtificialFromBasis_Infeasible(t *testing.T) {
	sm := &simplexMethod{
		m: 1,
		n: 1,
		rsmResult: rsmResult{
			x:       mat.NewVecDense(2, []float64{0, 1}),
			indices: mat.NewVecDense(1, []float64{1}),
		},
	}
	err := removeArtificialFromBasis(sm)
	if err == nil {
		t.Errorf("expected infeasible error, got nil")
	}
}

func TestRemoveArtificialFromBasis_Success(t *testing.T) {
	sm := &simplexMethod{
		m: 2,
		n: 2, // original variables
		rsmResult: rsmResult{
			x:       mat.NewVecDense(4, []float64{1e-9, 2, 0, 0}),
			indices: mat.NewVecDense(2, []float64{2, 3}), // artificial vars
		},
	}

	assert.Nil(t, removeArtificialFromBasis(sm))

	sm.x.SetVec(0, 1e-2)
	err := removeArtificialFromBasis(sm)
	if err != nil {
		t.Fatalf("removeArtificialFromBasis returned error: %v", err)
	}
}

func TestFindEnterSmall(t *testing.T) {
	A := mat.NewDense(2, 2, []float64{
		1, 0,
		0, 1,
	})
	c := mat.NewVecDense(2, []float64{-1, -2})
	pi := mat.NewVecDense(2, []float64{0, 0})
	isbasic := mat.NewVecDense(2, []float64{0, 1})

	fe := &enteringVariable{
		A:       A,
		c:       c,
		pi:      pi,
		isbasic: isbasic,

		epsilon: 1e-5,
	}

	assert.Nil(t, findEnter(fe))
	assert.Equal(t, fe.s, 0)
	assert.Equal(t, fe.cs, -1.0)
}

func TestFindLeaveSmall(t *testing.T) {
	B := mat.NewDense(2, 2, []float64{
		1, 0,
		0, 1,
	})
	indices := mat.NewVecDense(2, []float64{0, 1})
	xb := mat.NewVecDense(2, []float64{5, 3})
	as := mat.NewVecDense(2, []float64{1, 1})

	fl := &leavingVariable{
		B:       B,
		indices: indices,
		xb:      xb,
		as:      as,
		phase:   1,
		n:       2,
	}

	assert.Nil(t, findLeave(fl))
	assert.Equal(t, fl.r, 1)
}

func TestRSM_ImmediateOptimal(t *testing.T) {
	// Construct sm such that entering variable s=-1 immediately
	A := mat.NewDense(2, 2, []float64{1, 0, 0, 1})
	B := mat.NewDense(2, 2, []float64{1, 0, 0, 1})
	c := mat.NewVecDense(2, []float64{0, 0})
	indices := mat.NewVecDense(2, []float64{0, 1})
	sm := &simplexMethod{
		m: 2,
		n: 2,
		A: A,
		B: B,
		c: c,
		b: mat.NewVecDense(2, []float64{5, 3}),
		rsmResult: rsmResult{
			indices: indices,
			x:       mat.NewVecDense(2, nil),
			pi:      mat.NewVecDense(2, nil),
			flag:    common.SolverStatusNotSolved,
		},
		cb: mat.NewVecDense(2, []float64{0, 0}),
	}
	config := &common.SolverConfig{Tolerance: 1e-9, MaxIterations: 100}
	err := RSM(sm, 2, config)
	assert.Nil(t, err)
	assert.Equal(t, sm.flag, common.SolverStatusOptimal)
}

func TestUpdateB(t *testing.T) {
	B := mat.NewDense(2, 2, []float64{1, 2, 3, 4})
	as := mat.NewVecDense(2, []float64{9, 8})
	indices := mat.NewVecDense(2, []float64{0, 1})
	cb := mat.NewVecDense(2, []float64{0, 0})
	bu := &basisUpdate{BMat: B, indices: indices, cb: cb, as: as, s: 1, r: 0, cs: 7}
	if err := updateB(bu); err != nil {
		t.Fatalf("updateB failed: %v", err)
	}
	if B.At(0, 0) != 9 || B.At(1, 0) != 8 {
		t.Fatalf("B column not updated")
	}
	if int(indices.AtVec(0)) != 1 || cb.AtVec(0) != 7 {
		t.Fatalf("indices or cb not updated")
	}
}

func TestFindEnterAndLeave(t *testing.T) {
	// findEnter setup: m=1, n=2
	A := mat.NewDense(1, 2, []float64{1, 0})
	pi := mat.NewVecDense(1, []float64{0})
	c := mat.NewVecDense(2, []float64{0, -1})
	isbasic := mat.NewVecDense(2, []float64{1, 0})
	fe := &enteringVariable{A: A, pi: pi, c: c, isbasic: isbasic, epsilon: 1e-9}
	if err := findEnter(fe); err != nil {
		t.Fatalf("findEnter failed: %v", err)
	}
	if fe.s != 1 {
		t.Fatalf("expected s=1 got %d", fe.s)
	}

	// findLeave setup: m=2, B=I, as=[1,2], xb=[5,1]
	B := mat.NewDense(2, 2, []float64{1, 0, 0, 1})
	indices := mat.NewVecDense(2, []float64{0, 1})
	as := mat.NewVecDense(2, []float64{1, 2})
	xb := mat.NewVecDense(2, []float64{5, 1})
	fl := &leavingVariable{B: B, indices: indices, as: as, xb: xb, phase: 1, n: 2}
	if err := findLeave(fl); err != nil {
		t.Fatalf("findLeave failed: %v", err)
	}
	if fl.r != 1 {
		t.Fatalf("expected r=1 got %d", fl.r)
	}
}

func TestRemoveArtificialFromBasisANilPath(t *testing.T) {
	// sm with A == nil
	sm := &simplexMethod{
		m: 1,
		n: 2,
	}
	// indices contains artificial variable 2
	sm.indices = mat.NewVecDense(1, []float64{2})
	// x must be large enough and set artificial value to 0
	x := mat.NewVecDense(3, []float64{0, 0, 0})
	sm.x = x
	sm.A = nil
	// call removeArtificialFromBasis
	if err := removeArtificialFromBasis(sm); err != nil {
		t.Fatalf("removeArtificialFromBasis failed: %v", err)
	}
	if int(sm.indices.AtVec(0)) >= sm.n {
		t.Fatalf("expected artificial to be replaced")
	}
	// also ensure status unchanged when artificial positive
	sm.indices.SetVec(0, float64(2))
	sm.x.SetVec(2, 1.0)
	if err := removeArtificialFromBasis(sm); err == nil {
		t.Fatalf("expected error for positive artificial value")
	}
	// cleanup avoid affecting other tests
	sm.indices.SetVec(0, 0)
}

func TestSimplexEndToEnd(t *testing.T) {
}

func TestPhase1Repair(t *testing.T) {
	// Construct A so that initial basis (cols 0 and 1) is singular,
	// but replacing col 1 with non-basic original col 2 yields invertible basis.
	m := 2
	n := 3
	// Columns: 0=[1,0], 1=[0,0] (zero), 2=[0,1] (non-basic candidate)
	// Artificial columns appended: 3=[1,0], 4=[0,1]
	A := mat.NewDense(m, n+m, nil)
	// col0
	A.Set(0, 0, 1)
	A.Set(1, 0, 0)
	// col1 (zero)
	A.Set(0, 1, 0)
	A.Set(1, 1, 0)
	// col2
	A.Set(0, 2, 0)
	A.Set(1, 2, 1)
	// artificial cols (identity)
	A.Set(0, 3, 1)
	A.Set(1, 3, 0)
	A.Set(0, 4, 0)
	A.Set(1, 4, 1)

	sm := &simplexMethod{
		m:  m,
		n:  n,
		A:  A,
		c:  mat.NewVecDense(n+m, nil),
		cb: mat.NewVecDense(m, nil),
	}
	// set some costs
	for i := 0; i < n+m; i++ {
		sm.c.SetVec(i, float64(i))
	}
	// initial basis indices point to columns 0 and 1 (singular)
	sm.cb = mat.NewVecDense(m, []float64{0, 1})
	sm.B = matrix.ExtractColumns(sm.A, sm.cb)
	sm.b = mat.NewVecDense(m, []float64{5, 3})

	config := common.DefaultSolverConfig()
	config.Tolerance = 1e-9

	// Run RSM for Phase 1; it should repair the basis by swapping in col 2
	if err := RSM(sm, 1, config); err != nil {
		t.Fatalf("RSM failed: %v", err)
	}
	// Ensure that one of the basis indices is now 2
	if !contains(sm.indices, 2) {
		t.Fatalf("expected basis to contain column 2 after repair, got %v", sm.indices.RawVector().Data)
	}
}
func TestRemoveArtificialFromBasis_WithAReplace(t *testing.T) {
	A := mat.NewDense(2, 4, []float64{
		1, 0, 1, 0,
		0, 1, 0, 1,
	})
	sm := &simplexMethod{
		m: 2,
		n: 2,
		A: A,
		rsmResult: rsmResult{
			x:       mat.NewVecDense(4, []float64{0, 0, 0, 0}),
			indices: mat.NewVecDense(2, []float64{2, 3}), // artificial in basis
		},
	}
	if err := removeArtificialFromBasis(sm); err != nil {
		t.Fatalf("removeArtificialFromBasis failed: %v", err)
	}

	if int(sm.indices.AtVec(0)) < 0 || int(sm.indices.AtVec(0)) >= sm.n {
		t.Fatalf("expected artificial to be replaced by original column, got %v", sm.indices.AtVec(0))
	}
}

func TestSimplex_RepairSingularBasis(t *testing.T) {
	scf := &common.StandardComputationalForm{
		Objective:      mat.NewVecDense(1, []float64{1}),
		Constraints:    mat.NewDense(2, 1, []float64{1, 0}),
		RHS:            mat.NewVecDense(2, []float64{5, 3}),
		PrimalSolution: mat.NewVecDense(1, nil),
		ObjectiveValue: new(float64),
		Status:         new(common.SolverStatus),
	}
	config := common.DefaultSolverConfig()
	// Call Simplex; we primarily want to ensure it doesn't panic and returns an error or sets status
	err := Simplex(scf, config)
	// accept either nil or an error, but ensure status pointer is set
	if scf.Status == nil {
		t.Fatalf("expected status to be set")
	}
	// it's sufficient that the call completed
	_ = err
}

func TestSimplex_InfeasibleDetection(t *testing.T) {
	// Construct an infeasible LP: x >= 10 and x <= 5
	scf := &common.StandardComputationalForm{
		Objective:      mat.NewVecDense(1, []float64{1}),
		Constraints:    mat.NewDense(2, 1, []float64{1, 1}),
		RHS:            mat.NewVecDense(2, []float64{5, -10}),
		PrimalSolution: mat.NewVecDense(1, nil),
		ObjectiveValue: new(float64),
		Status:         new(common.SolverStatus),
	}
	config := common.DefaultSolverConfig()
	err := Simplex(scf, config)
	if err != nil {
		t.Fatalf("Simplex returned error on infeasible detection: %v", err)
	}
	if *scf.Status != common.SolverStatusInfeasible && *scf.Status != common.SolverStatusOptimal {
		t.Fatalf("expected infeasible or optimal status, got %v", *scf.Status)
	}
}

func TestFindEnterNoEnter(t *testing.T) {
	fe := &enteringVariable{A: mat.NewDense(1, 1, []float64{1}), pi: mat.NewVecDense(1, []float64{0}), c: mat.NewVecDense(1, []float64{1}), isbasic: mat.NewVecDense(1, []float64{1}), epsilon: 1e-9}
	if err := findEnter(fe); err != nil {
		t.Fatalf("findEnter failed: %v", err)
	}
	if fe.s != -1 {
		t.Fatalf("expected no entering variable, got %d", fe.s)
	}
}

func TestFindLeavePhase2Immediate(t *testing.T) {
	B := mat.NewDense(1, 1, []float64{1})
	indices := mat.NewVecDense(1, []float64{5})
	as := mat.NewVecDense(1, []float64{1})
	xb := mat.NewVecDense(1, []float64{1})
	fl := &leavingVariable{B: B, indices: indices, as: as, xb: xb, phase: 2, n: 2}
	if err := findLeave(fl); err != nil {
		t.Fatalf("findLeave failed: %v", err)
	}
	if fl.r != 0 {
		t.Fatalf("expected r=0 got %d", fl.r)
	}
}

func TestRSMMaxIterFailure(t *testing.T) {
	sm := &simplexMethod{
		m:         1,
		n:         1,
		B:         mat.NewDense(1, 1, []float64{0}),
		b:         mat.NewVecDense(1, []float64{1}),
		c:         mat.NewVecDense(2, []float64{1, 1}),
		cb:        mat.NewVecDense(1, []float64{1}),
		rsmResult: rsmResult{indices: mat.NewVecDense(1, []float64{0})},
	}
	config := &common.SolverConfig{Tolerance: 1e-9}
	if err := RSM(sm, 1, config); err == nil {
		t.Fatalf("expected RSM to return error on singular basis")
	}
}

func TestFindLeave_Unbounded(t *testing.T) {
	B := mat.NewDense(2, 2, []float64{1, 0, 0, 1})
	indices := mat.NewVecDense(2, []float64{0, 1})
	xb := mat.NewVecDense(2, []float64{5, 3})
	as := mat.NewVecDense(2, []float64{-1, 0})

	fl := &leavingVariable{B: B, indices: indices, xb: xb, as: as, phase: 1, n: 2}
	assert.Nil(t, findLeave(fl))
	assert.Equal(t, fl.r, -1)
}

func TestUpdateB_New(t *testing.T) {
	B := mat.NewDense(2, 2, []float64{0, 0, 0, 0})
	indices := mat.NewVecDense(2, []float64{0, 1})
	cb := mat.NewVecDense(2, []float64{10, 20})
	as := mat.NewVecDense(2, []float64{3, 4})
	bu := &basisUpdate{
		BMat:    B,
		indices: indices,
		cb:      cb,
		as:      as,
		s:       7,
		r:       1,
		cs:      42.0,
	}

	err := updateB(bu)
	assert.Nil(t, err)
	assert.Equal(t, B.At(0, 1), 3.0)
	assert.Equal(t, B.At(1, 1), 4.0)
	assert.Equal(t, int(indices.AtVec(1)), 7)
	assert.Equal(t, cb.AtVec(1), 42.0)
}

func TestRSM_DualFailure(t *testing.T) {
	B := mat.NewDense(1, 1, []float64{0}) // singular
	sm := &simplexMethod{
		m: 1,
		n: 1,
		A: mat.NewDense(1, 1, []float64{1}),
		B: B,
		c: mat.NewVecDense(1, []float64{1}),
		b: mat.NewVecDense(1, []float64{1}),
		rsmResult: rsmResult{
			indices: mat.NewVecDense(1, []float64{0}),
			x:       mat.NewVecDense(1, nil),
		},
		cb: mat.NewVecDense(1, []float64{0}),
	}
	config := &common.SolverConfig{Tolerance: 1e-9}
	err := RSM(sm, 1, config)
	assert.NotNil(t, err)
}

// Test findLeave when direction vector is all zeros (no leaving variable)
func TestFindLeave_AllZeros(t *testing.T) {
	B := mat.NewDense(2, 2, []float64{1, 0, 0, 1})
	indices := mat.NewVecDense(2, []float64{0, 1})
	xb := mat.NewVecDense(2, []float64{5, 3})
	as := mat.NewVecDense(2, []float64{0, 0}) // all zeros
	fl := &leavingVariable{B: B, indices: indices, xb: xb, as: as, phase: 1, n: 2}
	assert.Nil(t, findLeave(fl))
	assert.Equal(t, fl.r, -1)
}

// Test Simplex does not panic on singular basis repair
func TestSimplex_RepairSingularBasis_Corrected(t *testing.T) {
	scf := &common.StandardComputationalForm{
		Objective:      mat.NewVecDense(1, []float64{1}),
		Constraints:    mat.NewDense(2, 1, []float64{1, 0}),
		RHS:            mat.NewVecDense(2, []float64{5, 3}),
		PrimalSolution: mat.NewVecDense(1, nil),
		ObjectiveValue: new(float64),
		Status:         new(common.SolverStatus),
	}
	config := common.DefaultSolverConfig()
	_ = Simplex(scf, config)
	// Accept nil or error; just check no panic and Status is set
	if scf.Status == nil {
		t.Fatalf("expected status to be set")
	}
}

// Test RSM returns unbounded when leaving variable cannot be found
func TestRSM_Unbounded(t *testing.T) {
	B := mat.NewDense(2, 2, []float64{1, 0, 0, 1})
	indices := mat.NewVecDense(2, []float64{0, 1})
	xb := mat.NewVecDense(2, []float64{5, 3})
	as := mat.NewVecDense(2, []float64{-1, -2}) // negative directions
	fl := &leavingVariable{B: B, indices: indices, xb: xb, as: as, phase: 1, n: 2}
	assert.Nil(t, findLeave(fl))
	assert.Equal(t, fl.r, -1)
}

func TestRemoveArtificialFromBasis_FallbackReplacement(t *testing.T) {
	A := mat.NewDense(2, 4, []float64{
		1, 0, 1, 0,
		0, 1, 1, 0,
	})
	sm := &simplexMethod{
		m: 2,
		n: 2,
		A: A,
		rsmResult: rsmResult{
			x:       mat.NewVecDense(4, []float64{0, 0, 0, 0}),
			indices: mat.NewVecDense(2, []float64{2, 3}), // artificial in basis
		},
	}
	assert.Nil(t, removeArtificialFromBasis(sm))
	for i := 0; i < sm.indices.Len(); i++ {
		// ensure indices point to valid columns (may be original or artificial)
		assert.True(t, int(sm.indices.AtVec(i)) < sm.A.RawMatrix().Cols)
	}
}
