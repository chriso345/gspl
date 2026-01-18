package lp

import (
	"fmt"
	"math"
	"strings"
)

// String returns a string representation of the LinearProgram.
func (lp *LinearProgram) String() string {
	var sb strings.Builder

	sb.WriteString(lp.Description + "\n")

	if lp.Objective == nil {
		sb.WriteString("<undefined objective>\n")
		return sb.String()
	}

	if lp.Sense == LpMinimise {
		sb.WriteString("Minimize: ")
	} else {
		sb.WriteString("Maximize: ")
	}

	// Objective function
	first := true
	for i, coef := range lp.Objective.RawVector().Data {
		if coef == 0 {
			continue
		}
		if lp.Sense == LpMaximise {
			coef = -coef
		}
		if !first {
			if coef > 0 {
				sb.WriteString(" + ")
			} else {
				sb.WriteString(" - ")
			}
		} else {
			if coef < 0 {
				sb.WriteString("-")
			}
			first = false
		}
		varName := lp.Vars[i]
		fmt.Fprintf(&sb, "%.2f * %s", math.Abs(coef), varName.Name)
	}
	sb.WriteString("\n")

	if lp.Constraints == nil || lp.RHS == nil {
		sb.WriteString("<no constraints>\n")
		return sb.String()
	}

	// Constraints
	sb.WriteString("Subject to:\n")
	for row := range lp.Constraints.RawMatrix().Rows {
		fmt.Fprintf(&sb, "  C%d: ", row+1)

		first = true
		for col := range lp.Constraints.RawMatrix().Cols {
			coef := lp.Constraints.At(row, col)
			if lp.Vars[col].IsSlack {
				continue
			}
			if coef == 0 {
				continue
			}
			if !first {
				if coef > 0 {
					sb.WriteString(" + ")
				} else {
					sb.WriteString(" - ")
				}
			} else {
				if coef < 0 {
					sb.WriteString("-")
				}
				first = false
			}
			varName := lp.Vars[col]
			fmt.Fprintf(&sb, "%.2f * %s", math.Abs(coef), varName.Name)
		}

		// Constraint type and RHS
		switch lp.ConTypes[row] {
		case LpConstraintEQ:
			sb.WriteString(" == ")
		case LpConstraintGE:
			sb.WriteString(" >= ")
		case LpConstraintLE:
			sb.WriteString(" <= ")
		default:
			sb.WriteString(" ? ")
		}
		fmt.Fprintf(&sb, "%.3f\n", lp.RHS.AtVec(row))
	}

	// Variable bounds (integer, binary)
	intVars := []string{}
	binVars := []string{}
	for _, v := range lp.Vars {
		switch v.Category {
		case LpCategoryInteger:
			intVars = append(intVars, v.Name)
		case LpCategoryBinary:
			binVars = append(binVars, v.Name)
		}
	}

	if len(intVars) > 0 {
		sb.WriteString("Integer variables: " + strings.Join(intVars, ", ") + "\n")
	}
	if len(binVars) > 0 {
		sb.WriteString("Binary variables: " + strings.Join(binVars, ", ") + "\n")
	}

	return sb.String()
}

// PrintSolution prints the solution of the linear program in a human-readable format.
func (lp *LinearProgram) PrintSolution() {
	fmt.Println(lp.Description)
	// fmt.Printf("SolutionStatus: %s\n", lp.Status.String())
	fmt.Printf("ObjectiveValue: %.4f\n", lp.ObjectiveValue)

	for i, v := range lp.Vars {
		if !lp.Vars[i].IsSlack {
			fmt.Printf("%s: %.4f\n", v.Name, lp.PrimalSolution.AtVec(i))
		}
	}
}
