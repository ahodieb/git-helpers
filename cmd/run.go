package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func App() *cli.App {
	return &cli.App{
		Name:  AppName,
		Usage: "provides some additional functionality to git",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        FlagWorkingDir,
				Aliases:     []string{"C"},
				Usage:       "set working directory to `PATH`",
				DefaultText: ".",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "checkout-main",
				Usage:  "checkout `main` branch, it checks out the main branch by looking for `main` or `master`",
				Action: CheckoutMain,
			},
			{
				Name:   "rebase-all",
				Usage:  "rebase multiple branches onto `BRANCH`",
				Action: RebaseAll,
			},
			{
				Name:   "install-git-aliases",
				Usage:  "sets up git aliases for all helper subcommands",
				Action: InstallAliases,
			},
			{
				Name:   "version",
				Usage:  "shows the version of the cli",
				Action: Version,
			},
		},
	}
}

// Run executes the CLI with the given args and returns an exit code.
// Used by main and by the testscript harness so the CLI can be invoked
// in-process without calling os.Exit.
func Run(args []string) int {
	if err := App().Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
