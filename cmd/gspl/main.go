package main

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/chriso345/gspl/internal/lang"
	"github.com/chriso345/gspl/internal/lang/ast"
	"github.com/chriso345/gspl/solver"

	_ "github.com/chriso345/gspl/internal/lang/gmpl" // Required to register GMPL language
	_ "github.com/chriso345/gspl/internal/lang/mps"  // Required to register MPS language
)

var supportedLangs = []string{"gmpl", "mps"}

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

		// Get the language from the file extension
		ext := ""
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '.' {
				ext = path[i+1:]
				break
			}
		}

		if !slices.Contains(supportedLangs, ext) {
			exit(1, fmt.Errorf("unsupported file extension: %q", ext))
		}

		node, err := lang.ParseFile(ctx, ext, path)
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
