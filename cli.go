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

func (c *CLI) Run(args []string) int {
	var version bool
	flags := flag.NewFlagSet("insco", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.BoolVar(&version, "version", false, "Print version and quit")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParserFlagError
	}

	if version {
		fmt.Fprintf(c.errStream, "insco: version %s\n", Version)
		return ExitCodeOK
	}

	fmt.Fprintf(c.outStream, "insco works!\n")

	return ExitCodeOK
}
