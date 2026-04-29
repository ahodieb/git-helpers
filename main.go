package main

import (
	"os"

	"github.com/ahodieb/git-helpers/cmd"
)

func main() {
	os.Exit(cmd.Run(os.Args))
}
