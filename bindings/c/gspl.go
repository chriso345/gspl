package main

/*
#include "gspl.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

var lastError string

type program struct {
	lp   lp.LinearProgram
	vars []lp.LpVariable
	obj  lp.LpExpression
	objSense lp.LpSense
	constraints []struct{Expr lp.LpExpression; Type lp.LpConstraintType; RHS float64}
}

// --- program lifecycle ---

//export gspl_program_create
func gspl_program_create(name *C.char) C.GSPL_Handle {
	p := lp.NewLinearProgram(C.GoString(name), nil)
	return C.GSPL_Handle(put(&program{lp: p}))
}

//export gspl_program_free
func gspl_program_free(h C.GSPL_Handle) {
	del(handle(h))
}

// --- variables ---

//export gspl_program_add_variable
func gspl_program_add_variable(
	prog C.GSPL_Handle,
	name *C.char,
	cat C.GSPL_VarCategory,
) C.GSPL_Handle {

	p := get(handle(prog)).(*program)

	var category lp.LpCategory
	switch cat {
	case C.GSPL_INTEGER:
		category = lp.LpCategoryInteger
	case C.GSPL_BINARY:
		category = lp.LpCategoryBinary
	default:
		category = lp.LpCategoryContinuous
	}

	v := lp.NewVariable(C.GoString(name), category)
	p.vars = append(p.vars, v)
	p.lp.Vars = append(p.lp.Vars, v)

	return C.GSPL_Handle(put(v))
}

// --- objective ---

//export gspl_program_set_objective
func gspl_program_set_objective(
	prog C.GSPL_Handle,
	typ C.GSPL_ObjectiveType,
) {
	p := get(handle(prog)).(*program)
	mode := lp.LpMinimise
	if typ == C.GSPL_MAXIMIZE {
		mode = lp.LpMaximise
	}
	p.objSense = mode
	p.obj = lp.NewExpression(nil)
}

//export gspl_program_add_objective_term
func gspl_program_add_objective_term(
	prog C.GSPL_Handle,
	coeff C.double,
	varH C.GSPL_Handle,
) {
	p := get(handle(prog)).(*program)
	v := get(handle(varH)).(lp.LpVariable)
	p.obj.Terms = append(p.obj.Terms, lp.NewTerm(float64(coeff), v))
}

// --- constraints ---

//export gspl_program_add_constraint
func gspl_program_add_constraint(
	prog C.GSPL_Handle,
	typ C.GSPL_ConstraintType,
	rhs C.double,
) C.int {

	p := get(handle(prog)).(*program)

	ct := lp.LpConstraintLE
	switch typ {
	case C.GSPL_GE:
		ct = lp.LpConstraintGE
	case C.GSPL_EQ:
		ct = lp.LpConstraintEQ
	}

	id := len(p.constraints)
	p.constraints = append(p.constraints, struct{Expr lp.LpExpression; Type lp.LpConstraintType; RHS float64}{Expr: lp.NewExpression(nil), Type: ct, RHS: float64(rhs)})
	return C.int(id)
}

//export gspl_constraint_add_term
func gspl_constraint_add_term(
	prog C.GSPL_Handle,
	id C.int,
	coeff C.double,
	varH C.GSPL_Handle,
) {
	p := get(handle(prog)).(*program)
	v := get(handle(varH)).(lp.LpVariable)
	p.constraints[int(id)].Expr.Terms = append(p.constraints[int(id)].Expr.Terms, lp.NewTerm(float64(coeff), v))
}

// --- solving ---

//export gspl_program_solve
func gspl_program_solve(prog C.GSPL_Handle) C.GSPL_Handle {
	p := get(handle(prog)).(*program)
	// assemble LP from stored pieces
	p.lp.Vars = make([]lp.LpVariable, len(p.vars))
	copy(p.lp.Vars, p.vars)
	// set objective
	p.lp.AddObjective(p.objSense, p.obj)
	// add constraints
	for _, c := range p.constraints {
		p.lp.AddConstraint(c.Expr, c.Type, c.RHS)
	}

	sol, err := solver.Solve(&p.lp)
	if err != nil {
		lastError = err.Error()
		return 0
	}
	return C.GSPL_Handle(put(sol))
}

//export gspl_solution_objective_value
func gspl_solution_objective_value(sol C.GSPL_Handle) C.double {
	s := get(handle(sol)).(*solver.Solution)
	return C.double(s.ObjectiveValue)
}

//export gspl_solution_variable_count
func gspl_solution_variable_count(sol C.GSPL_Handle) C.int {
	s := get(handle(sol)).(*solver.Solution)
	return C.int(s.PrimalSolution.Len())
}

//export gspl_solution_variable_value
func gspl_solution_variable_value(sol C.GSPL_Handle, index C.int) C.double {
	s := get(handle(sol)).(*solver.Solution)
	return C.double(s.PrimalSolution.AtVec(int(index)))
}

//export gspl_solution_free
func gspl_solution_free(sol C.GSPL_Handle) {
	del(handle(sol))
}

// --- optional: strings ---

//export gspl_program_string
func gspl_program_string(prog C.GSPL_Handle) *C.char {
	p := get(handle(prog)).(*program)
	return C.CString(p.lp.String())
}

//export gspl_free_string
func gspl_free_string(str *C.char) {
	C.free(unsafe.Pointer(str))
}

// --- optional: last error ---

//export gspl_last_error
func gspl_last_error() *C.char {
	return C.CString(lastError)
}

// A no-op main to allow building as a C shared library.
func main() {}
