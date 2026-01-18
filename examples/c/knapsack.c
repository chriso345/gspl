#include <stdio.h>

// Generated C bindings for gspl
#include "../../bindings/c/libgspl.h"

int main() {
  // --- Knapsack problem ---
  GSPL_Handle knapsack = gspl_program_create("Knapsack Problem");

  // Values and weights
  double values[] = {5, 3, 6, 6, 2};
  double weights[] = {1, 4, 7, 6, 2};
  int n = 5;

  // Add binary variables x1..x5
  GSPL_Handle vars[5];
  for (int i = 0; i < n; i++) {
    char name[4];
    snprintf(name, sizeof(name), "x%d", i + 1);
    vars[i] = gspl_program_add_variable(knapsack, name, GSPL_BINARY);
  }

  // Objective: maximize total value
  gspl_program_set_objective(knapsack, GSPL_MAXIMIZE);
  for (int i = 0; i < n; i++) {
    gspl_program_add_objective_term(knapsack, values[i], vars[i]);
  }

  // Constraint: sum(weights[i] * x[i]) <= 15
  int c = gspl_program_add_constraint(knapsack, GSPL_LE, 15.0);
  for (int i = 0; i < n; i++) {
    gspl_constraint_add_term(knapsack, c, weights[i], vars[i]);
  }

  // Solve
  GSPL_Handle sol = gspl_program_solve(knapsack);
  printf("Knapsack Optimal Value: %f\n", gspl_solution_objective_value(sol));
  printf("Knapsack Solution: ");
  for (int i = 0; i < n; i++) {
    printf("%f ", gspl_solution_variable_value(sol, i));
  }
  printf("\n");

  // Cleanup
  gspl_solution_free(sol);
  gspl_program_free(knapsack);

  return 0;
}
