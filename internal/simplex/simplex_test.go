package simplex

import (
	"testing"

	"github.com/chriso345/gore/assert"
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
