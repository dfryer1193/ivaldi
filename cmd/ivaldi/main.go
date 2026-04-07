package main

import (
	"flag"
	"fmt"
	"os"

	"ivaldi/internal/cli"
)

func main() {
	fs := flag.NewFlagSet("ivaldi", flag.ExitOnError)

	var modulePrefix string
	fs.StringVar(&modulePrefix, "p", "", "Default module path prefix (e.g. github.com/user)")
	fs.StringVar(&modulePrefix, "module-prefix", "", "Default module path prefix")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <command>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "  init     Initialize a new Go project interactively\n")
		fmt.Fprintf(os.Stderr, "  update   Safely update existing Makefile and .golangci.yml\n")
		fmt.Fprintf(os.Stderr, "  clobber  Force overwrite Makefile, .golangci.yml, and CI\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		fs.PrintDefaults()
	}

	if len(os.Args) < 2 { //nolint:mnd // Minimal arguments check
		fs.Usage()
		os.Exit(1)
	}

	// The first argument is the command (init, update, clobber), following by flags
	// But standard flag package expects flags BEFORE positional args if using flag.Parse()
	// However, we want 'ivaldi [flags] <command>' or 'ivaldi <command> [flags]'?
	// The original code used flag.Parse() which expects flags before args.

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	args := fs.Args()
	if len(args) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	command := args[0]
	if command != "init" && command != "update" && command != "clobber" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fs.Usage()
		os.Exit(1)
	}

	if err := cli.Run(command, modulePrefix); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
