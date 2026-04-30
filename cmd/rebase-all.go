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
		log.Printf("%s branch has unpublished commits, aborting, %v", mainBranch, changes)
		return nil
	}

	branches, err := g.ListBranches()
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if branch == mainBranch {
			continue
		}
		if err := g.Rebase(mainBranch, branch); err != nil {
			log.Printf("failed to rebase %q: %v", branch, err)
		}
	}
	return nil
}
