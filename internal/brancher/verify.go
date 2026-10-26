package brancher

import (
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/internal/errors"
	"gonum.org/v1/gonum/mat"
)

// acceptIncumbent makes the candidate integer solution of scf the incumbent if
// it improves on it. The incumbent is kept in the original sense, and any
// improvement counts.
//
// Pure integer programs are checked exactly against the original program ip.SCF
// first (see verifyIntegerSolution), and a candidate that fails the check is an
// error rather than a solution. With continuous variables the candidate is only
// as exact as the float simplex.
func acceptIncumbent(ip *common.IntegerProgram, scf *common.StandardComputationalForm, config *common.SolverConfig) error {
	solution, obj := scf.PrimalSolution, *scf.ObjectiveValue
	if isPureInteger(ip.SCF) {
		var ok bool
		if solution, obj, ok = verifyIntegerSolution(ip.SCF, scf.PrimalSolution); !ok {
			return errors.New(errors.ErrNumericalFailure, "integer candidate is infeasible when checked exactly", nil)
		}
	}
	if scf.IsMaximization {
		obj = -obj
	}

	ip.BestMutex.Lock()
	defer ip.BestMutex.Unlock()
	if ip.BestSolution == nil ||
		(scf.IsMaximization && obj > ip.BestObj) ||
		(!scf.IsMaximization && obj < ip.BestObj) {
		ip.BestObj, ip.BestSolution = obj, solution
		if config.Debug {
			fmt.Printf("[DEBUG] New Best Obj: %.4f\n", ip.BestObj)
		}
	}
	return nil
}

// isPureInteger reports whether every non-slack column of scf is integer or binary.
func isPureInteger(scf *common.StandardComputationalForm) bool {
	for j, slackMarker := range scf.SlackIndices {
		if slackMarker != -1 {
			continue
		}
		if j >= len(scf.VarCategories) || scf.VarCategories[j] == common.VarCategoryContinuous {
			return false
		}
	}
	return true
}

// verifyIntegerSolution rounds the non-slack values of x and checks, in exact
// rational arithmetic, that the rounded point satisfies every row of scf. It
// returns the rounded point and its objective value in scf's (minimisation) sense,
// computed exactly and then rounded to float64, or ok=false if the point is
// infeasible.
//
// Coefficients are read as the shortest decimals that print as their float64
// values, so 0.1x + 0.2y <= 0.3 is checked as written although it fails for the
// float64 values. Rows must have the form lp produces: at most one slack or
// surplus column each.
func verifyIntegerSolution(scf *common.StandardComputationalForm, x *mat.VecDense) (*mat.VecDense, float64, bool) {
	m, n := scf.Constraints.Dims()
	point := mat.VecDenseCopyOf(x)
	values := make([]*big.Rat, n)
	obj := new(big.Rat)
	for j := range n {
		if scf.SlackIndices[j] != -1 {
			continue
		}
		v := math.Round(x.AtVec(j))
		if v < 0 || (scf.VarCategories[j] == common.VarCategoryBinary && v > 1) {
			return nil, 0, false
		}
		point.SetVec(j, v)
		values[j] = new(big.Rat).SetFloat64(v)

		c, ok := decimalRat(scf.Objective.AtVec(j))
		if !ok {
			return nil, 0, false
		}
		obj.Add(obj, c.Mul(c, values[j]))
	}

	for i := range m {
		// residual = b - A x over the non-slack columns
		residual, ok := decimalRat(scf.RHS.AtVec(i))
		if !ok {
			return nil, 0, false
		}
		slackCoef := 0.0
		slackCount := 0
		for j := range n {
			a := scf.Constraints.At(i, j)
			if a == 0 {
				continue
			}
			if scf.SlackIndices[j] != -1 {
				slackCoef = a
				slackCount++
				continue
			}
			aij, ok := decimalRat(a)
			if !ok {
				return nil, 0, false
			}
			residual.Sub(residual, aij.Mul(aij, values[j]))
		}

		// The slack takes up the residual and must be non-negative; without a
		// slack the row is an equality.
		switch slackCount {
		case 0:
			if residual.Sign() != 0 {
				return nil, 0, false
			}
		case 1:
			if float64(residual.Sign())*slackCoef < 0 {
				return nil, 0, false
			}
		default:
			return nil, 0, false
		}
	}
	objVal, _ := obj.Float64()
	return point, objVal, true
}

// decimalRat returns f as the shortest decimal that parses back to f, e.g. 1/10
// for 0.1. It fails for NaN and infinities.
func decimalRat(f float64) (*big.Rat, bool) {
	return new(big.Rat).SetString(strconv.FormatFloat(f, 'g', -1, 64))
}
