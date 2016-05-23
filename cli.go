package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"io"
	"os"
)

const (
	ExitCodeOK = iota
	ExitCodeParserFlagError
)

const (
	HomeEnv = "HOME"
)

type Options struct {
	OptHelp    bool `short:"h" long:"help" description:"Show this help message and exit"`
	OptVersion bool `short:"v" long:"version" description:"Print the version and exit"`
}

type CLI struct {
	outStream, errStream io.Writer
}

func (cli *CLI) showHelp() {
	// TODO: show usage
	fmt.Fprintf(cli.errStream, "Usage: \n")
}

func homeDir() string {
	homeDir := os.Getenv(HomeEnv)
	// TODO: error if HOME environment variablie does not exist
	return homeDir
}

func binaryDir() string {
	return homeDir() + "/bin"
}

func prefixDir() string {
	return homeDir() + "/usr/local"
}

func setup() error {
	binaryDir := binaryDir()

	if err := os.MkdirAll(binaryDir, os.ModePerm); err != nil {
		return err
	}

	prefixDir := prefixDir()
	if err := os.MkdirAll(prefixDir, os.ModePerm); err != nil {
		return err
	}

	return nil
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

	target := parsedArgs[0]
	// version := parsedArgs[1]

	setup()

	switch target {
	case "emacs":
	default:

	}

	return ExitCodeOK
}
