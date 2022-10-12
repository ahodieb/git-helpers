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

	mainBranch, err := g.FindMainBranch()
	if err != nil {
		return err
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
