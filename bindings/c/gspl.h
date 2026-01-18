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

/* --- optional: last error --- */
char *gspl_last_error(void);

#ifdef __cplusplus
}
#endif
