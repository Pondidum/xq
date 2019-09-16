package main

import (
	"fmt"
	"os"

	"xq/command"

	"github.com/mattn/go-colorable"
	"github.com/mitchellh/cli"
)

func main() {
	os.Exit(Run(os.Args[1:]))
}

func Run(args []string) int {

	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      colorable.NewColorableStdout(),
		ErrorWriter: colorable.NewColorableStderr(),
	}

	commands := command.Commands(ui)

	// deprecated command names
	hidden := []string{
		"helmRollback",
		"ReleaseHistory",
		"cleanDeployments",
		"cleanBranches",
		"netstorageCleanup",
		"dailyMerges",
	}

	cli := &cli.CLI{
		Name:                       "xq",
		Args:                       args,
		Commands:                   commands,
		HiddenCommands:             hidden,
		Autocomplete:               true,
		AutocompleteNoDefaultFlags: false,
		HelpFunc:                   cli.BasicHelpFunc("xq"),
		HelpWriter:                 os.Stdout,
	}

	exitCode, err := cli.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
