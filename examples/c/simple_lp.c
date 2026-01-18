/*
To build and run:
gcc examples/c/simple_lp.c -Ibindings/c -Lbindings/c -lgspl -o
examples/c/simple_lp LD_LIBRARY_PATH=bindings/c ./examples/c/simple_lp
*/

#include <stdio.h>

// Generated C bindings for gspl
#include "../../bindings/c/libgspl.h"

int main() {
  // Create a linear program
  GSPL_Handle prog = gspl_program_create("Example LP");

  // Add variables
  GSPL_Handle x1 = gspl_program_add_variable(prog, "x1", GSPL_CONTINUOUS);
  GSPL_Handle x2 = gspl_program_add_variable(prog, "x2", GSPL_CONTINUOUS);

  // Set objective (maximize x1 + 2*x2)
  gspl_program_set_objective(prog, GSPL_MAXIMIZE);
  gspl_program_add_objective_term(prog, 1.0, x1);
  gspl_program_add_objective_term(prog, 2.0, x2);

  // Add constraint: x1 + x2 <= 10
  int c1 = gspl_program_add_constraint(prog, GSPL_LE, 10.0);
  gspl_constraint_add_term(prog, c1, 1.0, x1);
  gspl_constraint_add_term(prog, c1, 1.0, x2);

  // Solve
  GSPL_Handle sol = gspl_program_solve(prog);
  printf("Optimal value: %f\n", gspl_solution_objective_value(sol));

  // Cleanup
  gspl_solution_free(sol);
  gspl_program_free(prog);
  return 0;
}
