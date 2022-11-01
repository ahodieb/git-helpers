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

	for _, command := range commands {
		err := g.AddAlias(command.Alias, fmt.Sprintf("!%s %s", AppName, command.Cmd))
		if err != nil {
			log.Printf("failed to add %q, %s", command.Cmd, err)
		}

		log.Printf("git %s - installed", command.Alias)
	}

	return nil
}
