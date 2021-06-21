package main

import (
	"fmt"
	"os"

	"github.com/jimeh/build-emacs-for-macos/pkg/cli"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	cliInstance := cli.New(version, commit, date)

	err := cliInstance.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		os.Exit(1)
	}
}
