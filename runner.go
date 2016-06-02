package main

import (
	"fmt"
	"io"
	"os/exec"
)

type Runner struct {
	outStream io.Writer
	err       error
}

func (runner *Runner) Err() error {
	return runner.err
}

func (runner *Runner) Run(cmd *exec.Cmd) {
	if runner.err != nil {
		return
	}
	out, err := cmd.Output()
	if err != nil {
		runner.err = err
		return
	}
	fmt.Fprintln(runner.outStream, out)
}
