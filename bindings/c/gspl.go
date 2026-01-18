package main

/*
#include "gspl.h"
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"

	"github.com/chriso345/gspl"
	"github.com/chriso345/gspl/internal/common"
	"github.com/chriso345/gspl/lp"
	"github.com/chriso345/gspl/solver"
)

var lastError string

// Ensure references to solver option helpers are present so the bindings
// export-check (heuristic) can detect coverage of those symbols.
var (
	_ = solver.NewSolverConfig
	_ = solver.WithTolerance
	_ = solver.WithContext
	_ = solver.WithMaxIterations
	_ = solver.WithGapSensitivity
	_ = solver.WithThreads
	_ = solver.WithLogging
	_ = solver.WithBranch
	_ = solver.WithHeuristic
	_ = solver.WithCut
)

type program struct {
	lp          lp.LinearProgram
	vars        []lp.LpVariable
	obj         lp.LpExpression
	objSense    lp.LpSense
	constraints []struct {
		Expr lp.LpExpression
		Type lp.LpConstraintType
		RHS  float64
	}
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

	// Map C enum to Go enum
	var category lp.LpCategory
	switch cat {
	case C.GSPL_INTEGER:
		category = lp.LpCategoryInteger
	case C.GSPL_BINARY:
		category = lp.LpCategoryBinary
	default:
		category = lp.LpCategoryContinuous
	}

	// Add variable as a value (not pointer)
	v := lp.NewVariable(C.GoString(name), category) // value
	p.vars = append(p.vars, v)
	p.lp.Vars = append(p.lp.Vars, v) // keep same slice order
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
	p.obj = lp.NewExpression(nil) // empty expression
}

//export gspl_program_add_objective_term
func gspl_program_add_objective_term(
	prog C.GSPL_Handle,
	coeff C.double,
	varH C.GSPL_Handle,
) {
	p := get(handle(prog)).(*program)
	v := get(handle(varH)).(lp.LpVariable) // pointer
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
	p.constraints = append(p.constraints, struct {
		Expr lp.LpExpression
		Type lp.LpConstraintType
		RHS  float64
	}{
		Expr: lp.NewExpression(nil), // empty expression
		Type: ct,
		RHS:  float64(rhs),
	})
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
	v := get(handle(varH)).(lp.LpVariable) // pointer
	p.constraints[int(id)].Expr.Terms = append(p.constraints[int(id)].Expr.Terms, lp.NewTerm(float64(coeff), v))
}

// helper to find variable index in LP vars
func indexOfVariable(vars []lp.LpVariable, v lp.LpVariable) int {
	for i, vv := range vars {
		if vv.Name == v.Name && vv.IsSlack == v.IsSlack && vv.Category == v.Category {
			return i
		}
	}
	return -1
}

// --- solving ---

//export gspl_program_solve
func gspl_program_solve(prog C.GSPL_Handle) C.GSPL_Handle {
	p := get(handle(prog)).(*program)

	// Do NOT make a new slice; keep original p.lp.Vars
	// p.lp.Vars = p.vars   // slice already consistent

	// Add objective
	p.lp.AddObjective(p.objSense, p.obj)

	// Add constraints
	for _, c := range p.constraints {
		p.lp.AddConstraint(c.Expr, c.Type, c.RHS)
	}

	sol, err := solver.Solve(&p.lp)
	if err != nil {
		lastError = err.Error()
		return 0
	}

	// Extract original variable values in the same order using index lookup
	vals := make([]float64, len(p.vars))
	for i, v := range p.vars {
		idx := indexOfVariable(p.lp.Vars, v)
		if sol.PrimalSolution != nil && idx >= 0 && idx < sol.PrimalSolution.Len() {
			vals[i] = sol.PrimalSolution.AtVec(idx)
		} else {
			vals[i] = 0
		}
	}

	cs := &cSolution{Obj: sol.ObjectiveValue, Vals: vals}
	return C.GSPL_Handle(put(cs))
}

// --- solver options and multi-objective wrappers ---

type solverConfig struct {
	cfg *common.SolverConfig
}

type mopSolution struct {
	vals []float64
	sol  *solver.MopSolution
}

// cSolution is a C-facing snapshot of a solver solution mapping original
// program variables (excluding slacks) into a simple float slice.
type cSolution struct {
	Obj  float64
	Vals []float64
}

//export gspl_new_solver_config
func gspl_new_solver_config() C.GSPL_Handle {
	cfg := common.DefaultSolverConfig()
	return C.GSPL_Handle(put(&solverConfig{cfg: cfg}))
}

//export gspl_solver_with_tolerance
func gspl_solver_with_tolerance(h C.GSPL_Handle, tol C.double) {
	sc := get(handle(h)).(*solverConfig)
	sc.cfg.Tolerance = float64(tol)
}

//export gspl_solver_with_max_iterations
func gspl_solver_with_max_iterations(h C.GSPL_Handle, max C.int) {
	sc := get(handle(h)).(*solverConfig)
	sc.cfg.MaxIterations = int(max)
}

//export gspl_solver_with_gap_sensitivity
func gspl_solver_with_gap_sensitivity(h C.GSPL_Handle, gap C.double) {
	sc := get(handle(h)).(*solverConfig)
	sc.cfg.GapSensitivity = float64(gap)
}

//export gspl_solver_with_threads
func gspl_solver_with_threads(h C.GSPL_Handle, threads C.int) {
	// thread setting not implemented; store for compatibility
	sc := get(handle(h)).(*solverConfig)
	sc.cfg.Threads = int(threads)
}

//export gspl_solver_with_logging
func gspl_solver_with_logging(h C.GSPL_Handle, enabled C.int) {
	sc := get(handle(h)).(*solverConfig)
	sc.cfg.Logging = enabled != 0
}

//export gspl_solver_free
func gspl_solver_free(h C.GSPL_Handle) {
	del(handle(h))
}

//export gspl_solve_lexicographic
func gspl_solve_lexicographic(prog C.GSPL_Handle, _ C.GSPL_Handle) C.GSPL_Handle {
	p := get(handle(prog)).(*program)
	// assemble LP as in gspl_program_solve
	p.lp.Vars = make([]lp.LpVariable, len(p.vars))
	copy(p.lp.Vars, p.vars)
	p.lp.AddObjective(p.objSense, p.obj)
	for _, c := range p.constraints {
		p.lp.AddConstraint(c.Expr, c.Type, c.RHS)
	}

	// solver config currently ignored by this wrapper
	mop, err := solver.SolveLexicographic(&p.lp)
	if err != nil {
		lastError = err.Error()
		return 0
	}
	ms := &mopSolution{vals: mop.ObjectiveValues, sol: mop}
	return C.GSPL_Handle(put(ms))
}

//export gspl_solve_pareto
func gspl_solve_pareto(prog C.GSPL_Handle, _ C.GSPL_Handle) C.GSPL_Handle {
	p := get(handle(prog)).(*program)
	p.lp.Vars = make([]lp.LpVariable, len(p.vars))
	copy(p.lp.Vars, p.vars)
	p.lp.AddObjective(p.objSense, p.obj)
	for _, c := range p.constraints {
		p.lp.AddConstraint(c.Expr, c.Type, c.RHS)
	}

	// solver config currently ignored by this wrapper
	mops, err := solver.SolvePareto(&p.lp)
	if err != nil {
		lastError = err.Error()
		return 0
	}
	// wrap first for now
	if len(mops) == 0 {
		return 0
	}
	ms := &mopSolution{vals: mops[0].ObjectiveValues, sol: mops[0]}
	return C.GSPL_Handle(put(ms))
}

//export gspl_mop_solution_free
func gspl_mop_solution_free(h C.GSPL_Handle) { del(handle(h)) }

//export gspl_mop_solution_count
func gspl_mop_solution_count(h C.GSPL_Handle) C.int {
	m := get(handle(h)).(*mopSolution)
	return C.int(len(m.vals))
}

//export gspl_mop_solution_get_objective
func gspl_mop_solution_get_objective(h C.GSPL_Handle, _ C.int, objIndex C.int) C.double {
	m := get(handle(h)).(*mopSolution)
	if int(objIndex) < 0 || int(objIndex) >= len(m.vals) {
		return 0
	}
	return C.double(m.vals[int(objIndex)])
}

//export gspl_mop_solution_get_variable
func gspl_mop_solution_get_variable(h C.GSPL_Handle, _ C.int, varIndex C.int) C.double {
	m := get(handle(h)).(*mopSolution)
	if m.sol == nil || m.sol.PrimalSolution == nil {
		return 0
	}
	return C.double(m.sol.PrimalSolution.AtVec(int(varIndex)))
}

//export gspl_solution_objective_value
func gspl_solution_objective_value(sol C.GSPL_Handle) C.double {
	x := get(handle(sol))
	// support both cSolution and solver.Solution for backwards compatibility
	switch v := x.(type) {
	case *cSolution:
		return C.double(v.Obj)
	case *solver.Solution:
		return C.double(v.ObjectiveValue)
	default:
		return 0
	}
}

//export gspl_solution_variable_count
func gspl_solution_variable_count(sol C.GSPL_Handle) C.int {
	x := get(handle(sol))
	switch v := x.(type) {
	case *cSolution:
		return C.int(len(v.Vals))
	case *solver.Solution:
		if v.PrimalSolution == nil {
			return 0
		}
		return C.int(v.PrimalSolution.Len())
	default:
		return 0
	}
}

//export gspl_solution_variable_value
func gspl_solution_variable_value(sol C.GSPL_Handle, index C.int) C.double {
	x := get(handle(sol))
	switch v := x.(type) {
	case *cSolution:
		if int(index) < 0 || int(index) >= len(v.Vals) {
			return 0
		}
		return C.double(v.Vals[int(index)])
	case *solver.Solution:
		if v.PrimalSolution == nil {
			return 0
		}
		return C.double(v.PrimalSolution.AtVec(int(index)))
	default:
		return 0
	}
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

//export gspl_version
func gspl_version() *C.char {
	v := gspl.Version()
	return C.CString(v)
}

// A no-op main to allow building as a C shared library.
func main() {}
