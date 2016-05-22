package main

import (
	"flag"
	"fmt"
	"io"
)

const (
	ExitCodeOK = iota
	ExitCodeParserFlagError
)

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI) Run(args []string) int {
	var version bool
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&version, "version", false, "display the version")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParserFlagError
	}

	if version {
		fmt.Fprintf(cli.errStream, "%s: v%s\n", Name, Version)
		return ExitCodeOK
	}

	parsedArgs := flags.Args()

	if len(parsedArgs) == 0 {
		fmt.Fprintf(cli.errStream, "[Error]: You must specify the target.\n")
		// TODO: show usage
		return ExitCodeParserFlagError
	}

	fmt.Fprintf(cli.outStream, "insco works!\n")

	return ExitCodeOK
}
