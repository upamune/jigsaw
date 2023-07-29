package main

import (
	"fmt"
	"os"
)

const (
	ExitOK int = 0
	ExitNG int = 1
)

var (
	Version  string
	Revision string
)

func main() {
	os.Exit(run())
}

func run() int {
	cli := &CLI{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Stdin:  os.Stdin,
	}
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(cli.Stderr, "ERROR: %s\nVersion: %s\nRevision: %s\n", err.Error(), Version, Revision)
		return ExitNG
	}
	return ExitOK
}
