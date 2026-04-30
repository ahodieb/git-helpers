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
	onto := fs.String("onto", "", "ref to rebase onto; when unset, falls back to local main or master")
	if err := fs.Parse(args[1:]); err != nil {
		return 2
	}

	if err := rebaseAll(git.Git{WorkingDir: *dir}, *onto); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func rebaseAll(g git.Git, onto string) error {
	base, err := resolveBase(g, onto)
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

func resolveBase(g git.Git, onto string) (string, error) {
	if onto != "" {
		return onto, nil
	}
	for _, ref := range []string{"main", "master"} {
		if _, err := g.Exec("rev-parse", "--verify", ref); err == nil {
			return ref, nil
		}
	}
	return "", fmt.Errorf("no main or master found; pass --onto <ref>")
}
