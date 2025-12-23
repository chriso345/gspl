package brancher

import (
	"testing"

	"github.com/chriso345/gore/assert"
	"github.com/chriso345/gspl/internal/common"
	"gonum.org/v1/gonum/mat"
)

func TestDefaultBranch_HappyPath(t *testing.T) {
	scf := &common.StandardComputationalForm{
		PrimalSolution: mat.NewVecDense(3, []float64{1.0, 2.5, 3.0}),
		Constraints:    mat.NewDense(1, 3, []float64{0, 0, 0}),
		RHS:            mat.NewVecDense(1, []float64{0}),
		Objective:      mat.NewVecDense(3, []float64{0, 0, 0}),
		ObjectiveValue: new(float64),
		Status:         new(common.SolverStatus),
		SlackIndices:   make([]int, 3),
	}
	node := &common.Node{SCF: scf}

	children, err := DefaultBranch(node)
	assert.Nil(t, err)
	assert.Equal(t, len(children), 2)

	up := children[0]
	down := children[1]

	r := down.SCF.Constraints.RawMatrix().Rows - 1
	val := down.SCF.Constraints.At(r, 1)
	assert.Equal(t, val, 1.0)
	rhs := down.SCF.RHS.AtVec(r)
	assert.Equal(t, rhs, 2.0)

	// Check up (dir 2) has a -1 at column 1 and RHS -3 (-(int+1))
	r2 := up.SCF.Constraints.RawMatrix().Rows - 1
	val2 := up.SCF.Constraints.At(r2, 1)
	assert.Equal(t, val2, -1.0)
	rhs2 := up.SCF.RHS.AtVec(r2)
	assert.Equal(t, rhs2, -3.0)
}

func TestDefineStrategies_SetsDefaultsOrUsesProvided(t *testing.T) {
	ip := &common.IntegerProgram{}
	defineStrategies(ip)
	if branchFunc == nil {
		t.Fatalf("expected branchFunc to be set")
	}

	called := false
	myBranch := func(n *common.Node) ([]*common.Node, error) { called = true; return nil, nil }
	ip2 := &common.IntegerProgram{Branch: myBranch}
	defineStrategies(ip2)
	scf := &common.StandardComputationalForm{
		PrimalSolution: mat.NewVecDense(1, []float64{0}),
		Constraints:    mat.NewDense(1, 1, []float64{0}),
		RHS:            mat.NewVecDense(1, []float64{0}),
	}
	_, err := branchFunc(&common.Node{SCF: scf})
	if err != nil {
		t.Fatalf("branchFunc returned error: %v", err)
	}
	if !called {
		t.Fatalf("expected branchFunc to call provided function")
	}
}

func TestDefaultHeuristicAndCut(t *testing.T) {
	vals, obj, ok := DefaultHeuristic(nil)
	assert.Equal(t, obj, 0.0)
	assert.Equal(t, ok, false)
	if vals != nil {
		assert.Equal(t, len(vals), 0)
	}

	cuts := DefaultCut(nil)
	if cuts != nil {
		assert.Equal(t, len(cuts), 0)
	}
}
