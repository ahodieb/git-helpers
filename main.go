package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ahodieb/git-helpers/git"
)

func main() {
	os.Exit(run(os.Args))
}

func run(args []string) int {
	fs := flag.NewFlagSet("git-helpers", flag.ContinueOnError)
	dir := fs.String("C", ".", "working directory")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}

	if err := rebaseAll(git.Git{WorkingDir: *dir}); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// rebaseAll replays every local branch (other than main/master) onto
// main using `git replay`. Branches that conflict are reported and skipped.
func rebaseAll(g git.Git) error {
	base, err := g.FindMainBranch()
	if err != nil {
		return err
	}

	branches, err := g.ListBranches()
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if branch == base {
			continue
		}
		if err := g.Replay(base, branch); err != nil {
			log.Printf("failed to replay %q: %v", branch, err)
		}
	}
	return nil
}
