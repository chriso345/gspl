package main

import (
	"context"
	"fmt"
	"os"

	"github.com/chriso345/gspl/internal/lang"
	"github.com/chriso345/gspl/internal/lang/ast"
	_ "github.com/chriso345/gspl/internal/lang/gmpl"
	"github.com/chriso345/gspl/solver"
)

func exit(code int, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	os.Exit(code)
}

func main() {
	args := ParseArgs()

	if args.Run.Subcommand {
		path := args.Run.File.Value
		fmt.Println("Running file:", path)
		ctx := context.Background()
		node, err := lang.ParseFile(ctx, "gmpl", path)
		if err != nil {
			exit(1, err)
		}
		fmt.Printf("Parsed node: %T\n", node)
		if m, ok := node.(*ast.Module); ok && m.LP != nil {
			fmt.Println("Found linear program; solving...")
			sol, err := solver.Solve(m.LP)
			if err != nil {
				exit(1, err)
			}
			fmt.Printf("Status: %v\n", sol.Status)
			fmt.Printf("Objective: %.6f\n", sol.ObjectiveValue)
			fmt.Printf("Primal: %v\n", sol.PrimalSolution.RawVector().Data)
		}
	}

	exit(0, nil)
}
