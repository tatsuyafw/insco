package main

import "os"

const Name string = "insco"
const Version string = "0.0.1"

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args[1:]))
}
