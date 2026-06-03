package main

import (
	"os"

	"fontview/cmd"
)

var Version = "0.0.0-dev"

func main() {
	if err := cmd.Execute(os.Args[1:], Version); err != nil {
		cmd.PrintError(err)
		os.Exit(1)
	}
}
