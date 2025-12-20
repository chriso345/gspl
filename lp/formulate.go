package lp

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
)

// AddObjective adds a new objective function to the LP model.
func (lp *LinearProgram) AddObjective(sense LpSense, expr LpExpression) {
	lp.Sense = sense

	// Create the objective function vector
	lp.Objective = mat.NewVecDense(len(lp.Vars), nil)

	for _, term := range expr.Terms {
		// Find the index of the variable in lp.Vars
		varIndex := -1
		for i, v := range lp.Vars {
			if v.Name == term.Variable.Name {
				varIndex = i
				break
			}
		}
		if varIndex == -1 {
			panic(fmt.Sprintf("variable in objective function not found in LP variables: %s", term.Variable.Name))
		}
		coef := term.Coefficient

		// Solver requires a minimisation problem. If maximising we invert the coefficients
		// here and mark the LinearProgram so newSCF does not invert a second time.
		if sense == LpMaximise {
			coef = -coef
			lp.ObjectiveIsNegated = true
		}

		lp.Objective.SetVec(varIndex, coef)
	}
}

// AddConstraint adds a new constraint to the LP model.
func (lp *LinearProgram) AddConstraint(expr LpExpression, conType LpConstraintType, rhs float64) {
	if lp.Objective == nil {
		panic("objective function must be defined before adding constraints")
	}

	// Create a new row for the constraint
	newRow := make([]float64, len(lp.Vars))

	// If rhs < 0, invert the constraint
	flipped := false
	if rhs < 0 {
		rhs = -rhs
		flipped = true
		switch conType {
		case LpConstraintLE:
			conType = LpConstraintGE
		case LpConstraintGE:
			conType = LpConstraintLE
		}
	}

	for _, term := range expr.Terms {
		// Find the index of the variable in lp.Vars
		varIndex := -1
		for i, v := range lp.Vars {
			if v.Name == term.Variable.Name {
				varIndex = i
				break
			}
		}
		if varIndex == -1 {
			panic(fmt.Sprintf("variable in constraint not found in LP variables: %s", term.Variable.Name))
		}

		if flipped {
			newRow[varIndex] = -term.Coefficient
		} else {
			newRow[varIndex] = term.Coefficient
		}
	}

	// Append the new constraint to the relevant matrices
	if lp.Constraints == nil {
		lp.Constraints = mat.NewDense(1, len(lp.Vars), newRow)
		lp.RHS = mat.NewVecDense(1, []float64{rhs})
		lp.ConTypes = []LpConstraintType{conType}
	} else {
		lp.Constraints = mat.NewDense(lp.Constraints.RawMatrix().Rows+1, len(lp.Vars), append(lp.Constraints.RawMatrix().Data, newRow...))
		lp.RHS = mat.NewVecDense(lp.RHS.Len()+1, append(lp.RHS.RawVector().Data, rhs))
		lp.ConTypes = append(lp.ConTypes, conType)
	}

	// If the constraint is GE or LE we need to add a slack/surplus variable
	if conType == LpConstraintLE || conType == LpConstraintGE {
		slackVar := LpVariable{
			Name:     "s" + fmt.Sprint(len(lp.Vars)+1),
			IsSlack:  true,
			Category: LpCategoryContinuous,
		}
		lp.Vars = append(lp.Vars, slackVar)

		// Expand the objective function to include the new slack variable with coefficient 0
		newObjective := mat.NewVecDense(lp.Objective.Len()+1, nil)
		for i := 0; i < lp.Objective.Len(); i++ {
			newObjective.SetVec(i, lp.Objective.AtVec(i))
		}
		newObjective.SetVec(newObjective.Len()-1, 0) // Coefficient of slack variable is 0
		lp.Objective = newObjective

		currentRowIndex := lp.Constraints.RawMatrix().Rows - 1

		// Add a single new column for this constraint
		newConstraints := mat.NewDense(
			lp.Constraints.RawMatrix().Rows,
			lp.Constraints.RawMatrix().Cols+1,
			nil,
		)

		// Copy old values
		for r := 0; r < lp.Constraints.RawMatrix().Rows; r++ {
			for c := 0; c < lp.Constraints.RawMatrix().Cols; c++ {
				newConstraints.Set(r, c, lp.Constraints.At(r, c))
			}
		}

		// Set the slack/surplus for only the current row
		switch conType {
		case LpConstraintLE:
			newConstraints.Set(currentRowIndex, newConstraints.RawMatrix().Cols-1, 1) // slack
		case LpConstraintGE:
			newConstraints.Set(currentRowIndex, newConstraints.RawMatrix().Cols-1, -1) // surplus
		}

		lp.Constraints = newConstraints
	}
}
