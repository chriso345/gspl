package main

import (
	"os"
	"testing"
)

func TestParseArgs_Default(t *testing.T) {
	orig := os.Args
	defer func() { os.Args = orig }()

	os.Args = []string{"gspl"}
	args := ParseArgs()
	if args == nil {
		t.Fatal("expected non-nil args")
	}
	if args.Version.Subcommand {
		t.Fatal("version should not be set for default args")
	}
	if args.Run.Subcommand {
		t.Fatal("run should not be set for default args")
	}
}

func TestParseArgs_Run(t *testing.T) {
	orig := os.Args
	defer func() { os.Args = orig }()

	// Try a couple of common flag syntaxes; clifford may accept -file or --file or -file=/path
	os.Args = []string{"gspl", "run", "file.txt"}
	args := ParseArgs()
	if args == nil {
		t.Fatal("expected non-nil args")
	}
	if !args.Run.Subcommand {
		t.Fatal("expected run subcommand to be set")
	}

	if args.Run.File.Value != "file.txt" {
		t.Fatalf("unexpected file value: %q", args.Run.File.Value)
	}
}
