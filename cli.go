package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
)

const (
	ExitCodeOK = iota
	ExitCodeParserFlagError
)

func (cli *CLI) showHelp() {
	// TODO: show usage
	fmt.Fprintf(cli.errStream, "Usage: \n")
}

type CLI struct {
	outStream, errStream io.Writer
}

type Options struct {
	OptHelp    bool `short:"h" long:"help" description:"Show this help message and exit"`
	OptVersion bool `short:"v" long:"version" description:"Print the version and exit"`
}

func (cli *CLI) parseOptions(args []string) (*Options, []string, error) {
	opts := &Options{}
	p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.ParseArgs(args)
	if err != nil {
		return nil, nil, err
	}

	return opts, args, nil
}

func (cli *CLI) Run(args []string) int {
	opts, parsedArgs, err := cli.parseOptions(args)
	if err != nil {
		return ExitCodeParserFlagError
	}

	if opts.OptHelp {
		cli.showHelp()
		return ExitCodeOK
	}

	if opts.OptVersion {
		fmt.Fprintf(cli.errStream, "%s: v%s\n", Name, Version)
		return ExitCodeOK
	}

	if len(parsedArgs) == 0 {
		fmt.Fprintf(cli.errStream, "[Error]: You must specify the target.\n")
		cli.showHelp()
		return ExitCodeParserFlagError
	}

	fmt.Println(parsedArgs)

	return ExitCodeOK
}
