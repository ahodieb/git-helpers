package main

import (
	"log"
	"os"

	"github.com/ahodieb/git-helpers/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  cmd.AppName,
		Usage: "provides some additional functionality to git",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        cmd.FlagWorkingDir,
				Aliases:     []string{"C"},
				Usage:       "set working directory to `PATH`",
				DefaultText: ".",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "checkout-main",
				Usage:  "checkout `main` branch, it checks out the main branch by looking for `main` or `master`",
				Action: cmd.CheckoutMain,
			},
			{
				Name:   "rebase-all",
				Usage:  "rebase multiple branches onto `BRANCH`",
				Action: cmd.RebaseAll,
			},
			{
				Name:   "install-git-aliases",
				Usage:  "sets up git aliases for all helper subcommands",
				Action: cmd.InstallAliases,
			},
			{
				Name:   "version",
				Usage:  "shows the version of the cli",
				Action: cmd.Version,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
