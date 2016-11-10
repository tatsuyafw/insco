package main

import (
	"os/exec"
	"testing"
)

type fakeWriter struct{}

func (f fakeWriter) Write([]byte) (n int, e error) {
	return
}

func TestNewRunner(t *testing.T) {
	runner := &Runner{outStream: fakeWriter{}}
	if runner.Err() != nil {
		t.Errorf("New Runner has a not nil error")
	}
}

func TestRunner(t *testing.T) {
	runner := &Runner{outStream: fakeWriter{}}
	runner.Run(exec.Command("dir"))
	if runner.Err() != nil {
		t.Errorf("New Runner has a not nil error")
	}
}

func TestRunnerError1(t *testing.T) {
	runner := &Runner{outStream: fakeWriter{}}
	runner.Run(exec.Command("no_such_command"))
	if runner.Err() == nil {
		t.Errorf("It is expected that Runner has an error, but does not have.")
	}
}

func TestRunnerError2(t *testing.T) {
	runner := &Runner{outStream: fakeWriter{}}
	runner.Run(exec.Command("dir"))
	runner.Run(exec.Command("no_such_command"))
	if runner.Err() == nil {
		t.Errorf("It is expected that Runner has an error, but does not have.")
	}
}
