package cmd

import (
	"github.com/ahodieb/git-helpers/git"
	"github.com/urfave/cli/v2"
)

func CheckoutMain(cCtx *cli.Context) error {
	g := git.Git{WorkingDir: cCtx.String(FlagWorkingDir)}

	mainBranch, err := g.FindMainBranch()
	if err != nil {
		return err
	}

	return g.Checkout(mainBranch)
}
