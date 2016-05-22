package main

import (
	"fmt"
	"os"
)

const Name string = "insco"
const Version string = "0.0.1"

func main() {
	fmt.Println("Hello, insco!")

	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args[1:]))
}
