package main

import (
	"flag"
	"fmt"
	"os"

	"ivaldi/internal/cli"
)

func main() {
	var modulePrefix string
	flag.StringVar(&modulePrefix, "p", "", "Default module path prefix (e.g. github.com/user)")
	flag.StringVar(&modulePrefix, "module-prefix", "", "Default module path prefix")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <command>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nCommands:\n")
		fmt.Fprintf(os.Stderr, "  init     Initialize a new Go project interactively\n")
		fmt.Fprintf(os.Stderr, "  update   Safely update existing Makefile and .golangci.yml\n")
		fmt.Fprintf(os.Stderr, "  clobber  Force overwrite Makefile, .golangci.yml, and CI\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	command := args[0]
	if command != "init" && command != "update" && command != "clobber" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		flag.Usage()
		os.Exit(1)
	}

	if err := cli.Run(command, modulePrefix); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
