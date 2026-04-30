package cmd

import (
	"fmt"
	"log"

	"github.com/ahodieb/git-helpers/git"
	"github.com/urfave/cli/v2"
)

var commands = []struct {
	Alias string
	Cmd   string
}{
	{Alias: "main", Cmd: "checkout-main"},
	{Alias: "rebase-all", Cmd: "rebase-all"},
}

func InstallAliases(cCtx *cli.Context) error {
	g := git.Git{WorkingDir: cCtx.String(FlagWorkingDir)}

	for _, c := range commands {
		if err := g.AddAlias(c.Alias, fmt.Sprintf("!%s %s", AppName, c.Cmd)); err != nil {
			log.Printf("failed to add %q: %v", c.Cmd, err)
			continue
		}
		log.Printf("git %s - installed", c.Alias)
	}
	return nil
}
