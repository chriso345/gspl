#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

/* --- handles --- */
typedef uint64_t GSPL_Handle;

/* --- enums --- */
typedef enum {
  GSPL_CONTINUOUS = 0,
  GSPL_INTEGER = 1,
  GSPL_BINARY = 2
} GSPL_VarCategory;

typedef enum { GSPL_MINIMIZE = 0, GSPL_MAXIMIZE = 1 } GSPL_ObjectiveType;

typedef enum { GSPL_LE = 0, GSPL_GE = 1, GSPL_EQ = 2 } GSPL_ConstraintType;

/* --- program lifecycle --- */
GSPL_Handle gspl_program_create(char *name);
void gspl_program_free(GSPL_Handle prog);

/* --- variables --- */
GSPL_Handle gspl_program_add_variable(GSPL_Handle prog, char *name,
                                      GSPL_VarCategory cat);

/* --- objective --- */
void gspl_program_set_objective(GSPL_Handle prog, GSPL_ObjectiveType type);

void gspl_program_add_objective_term(GSPL_Handle prog, double coeff,
                                     GSPL_Handle var);

/* --- constraints --- */
int gspl_program_add_constraint(GSPL_Handle prog, GSPL_ConstraintType type,
                                double rhs);

void gspl_constraint_add_term(GSPL_Handle prog, int constraint_id, double coeff,
                              GSPL_Handle var);

/* --- solving --- */
GSPL_Handle gspl_program_solve(GSPL_Handle prog);

/* --- solution access --- */
double gspl_solution_objective_value(GSPL_Handle sol);
int gspl_solution_variable_count(GSPL_Handle sol);
double gspl_solution_variable_value(GSPL_Handle sol, int index);
void gspl_solution_free(GSPL_Handle sol);

/* --- optional: strings --- */
char *gspl_program_string(GSPL_Handle prog);
void gspl_free_string(char *str);

/* --- multi-objective solving --- */
GSPL_Handle gspl_solve_lexicographic(GSPL_Handle prog, GSPL_Handle solver_cfg);
GSPL_Handle gspl_solve_pareto(GSPL_Handle prog, GSPL_Handle solver_cfg);
void gspl_mop_solution_free(GSPL_Handle mop);
int gspl_mop_solution_count(GSPL_Handle mop);
double gspl_mop_solution_get_objective(GSPL_Handle mop, int solIndex, int objIndex);
double gspl_mop_solution_get_variable(GSPL_Handle mop, int solIndex, int varIndex);

/* --- solver options --- */
GSPL_Handle gspl_new_solver_config(void);
void gspl_solver_with_tolerance(GSPL_Handle cfg, double tol);
void gspl_solver_with_max_iterations(GSPL_Handle cfg, int max);
void gspl_solver_with_gap_sensitivity(GSPL_Handle cfg, double gap);
void gspl_solver_with_threads(GSPL_Handle cfg, int threads);
void gspl_solver_with_logging(GSPL_Handle cfg, int enabled);
void gspl_solver_free(GSPL_Handle cfg);

/* --- optional: strings --- */
char *gspl_program_string(GSPL_Handle prog);
void gspl_free_string(char *str);

/* --- optional: last error --- */
char *gspl_last_error(void);

/* --- version --- */
char *gspl_version(void);

#ifdef __cplusplus
}
#endif
