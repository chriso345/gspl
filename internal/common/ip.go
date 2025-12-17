package common

import (
	"gonum.org/v1/gonum/mat"
	"sync"
)

type IntegerProgram struct {
	SCF *StandardComputationalForm

	// IP Specific fields
	ConstraintDir []string

	// Best known solution
	BestSolution *mat.VecDense // x*
	BestObj      float64
	// Mutex to protect BestObj and BestSolution updates across goroutines
	BestMutex sync.Mutex

	// User-supplied strategy functions
	Branch    BranchFunc
	Heuristic HeuristicFunc
	Cut       CutFunc
}

// FIXME: This is just a placeholder struct for Node. This will change.
type Node struct {
	SCF *StandardComputationalForm

	ID       int
	ParentID int
	Depth    int

	Bounds [][2]float64

	RelaxedSol []float64
	RelaxedObj float64

	IsFeasible bool
	IsInteger  bool

	BranchVar   int
	BranchValue float64

	LowerBound float64
}
