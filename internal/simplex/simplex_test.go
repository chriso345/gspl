package simplex

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/internal/common"
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
	config := &common.SolverConfig{Tolerance: 1e-9}
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
	// exercise RSM early failure by creating singular B when solving xb
	// create sm with m=1 n=1, B zero matrix will cause SolveVec to error
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
