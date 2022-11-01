package main

import (
	"log"
	"os"

	"github.com/ahodieb/rebaser/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "git-helpers",
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
				Name:   "rebase-all",
				Usage:  "rebase multiple branches onto `BRANCH`",
				Action: cmd.RebaseAll,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}