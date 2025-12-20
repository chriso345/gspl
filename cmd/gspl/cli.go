package main

import (
	"fmt"
	"os"

	"github.com/chriso345/clifford"
	"github.com/chriso345/gspl"
)

type CLIArgs struct {
	clifford.Clifford `name:"gspl"`
	clifford.Help     `type:"subcmd"`

	Run struct {
		clifford.Subcommand
		clifford.Desc `desc:"Run a specific linear program"`

		File struct {
			Value string
			clifford.Required
			clifford.Desc `desc:"Path to the linear program file to run"`
		}
	}

	Version struct {
		clifford.Subcommand
		clifford.Desc `desc:"Show version information"`
	}
}

func ParseArgs() *CLIArgs {
	args := &CLIArgs{}

	if err := clifford.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing arguments:", err)
		os.Exit(1)
	}

	if args.Version.Subcommand {
		version := gspl.Version()
		fmt.Printf("gspl version %s\n", version)
		os.Exit(0)
	}

	return args
}
