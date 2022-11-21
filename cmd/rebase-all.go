package cmd

import (
	"log"

	"github.com/ahodieb/git-helpers/git"
	"github.com/urfave/cli/v2"
)

func RebaseAll(cCtx *cli.Context) error {
	g := git.Git{WorkingDir: cCtx.String(FlagWorkingDir)}

	mainBranch, err := g.FindMainBranch()
	if err != nil {
		return err
	}

	changes, err := g.Cherry()
	if err != nil {
		return err
	}

	if len(changes) != 0 {
		log.Printf("%s branch has unpublished commits, abborting", mainBranch)
		return nil
	}

	branches, err := g.ListBranches()
	if err != nil {
		return err
	}

	branches = remove(mainBranch, branches)

	for _, branch := range branches {
		if err := g.Rebase(mainBranch, branch); err != nil {
			log.Printf("failed to rebase branch %q", branch)
		}
	}

	return err
}

func remove(value string, values []string) []string {
	for i := range values {
		if values[i] == value {
			return append(values[:i], values[i+1:]...)
		}
	}

	return values
}
