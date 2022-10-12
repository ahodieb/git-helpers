package main

import (
	"log"
	"os"

	"github.com/ahodieb/rebaser/git"
	"github.com/urfave/cli/v2"
)

const FlagWorkingDir = "working-dir"

func main() {
	app := &cli.App{
		Name:   "rebaser",
		Usage:  "rebase multiple branches onto main or master",
		Action: RebaseAllBranches,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        FlagWorkingDir,
				Aliases:     []string{"C"},
				Usage:       "set working directory `PATH`",
				DefaultText: ".",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func RebaseAllBranches(cCtx *cli.Context) error {
	g := git.Git{WorkingDir: cCtx.String(FlagWorkingDir)}
	stdout, stderr, err := g.Exec("status")

	log.Println(stdout)
	log.Println(stderr)
	return err
}
