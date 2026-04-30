package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"git-helpers": func() int { return run(os.Args) },
	}))
}

// initRepo runs in each scenario's WorkDir before the script starts.
// Creates ./repo with one initial commit on `main`, configured so
// commits don't depend on a global git identity or signing key.
const initRepo = `
set -e
git init -q -b main repo
cd repo
git config user.email 'test@git-helpers.local'
git config user.name 'git-helpers tests'
git config commit.gpgsign false
git commit -q --allow-empty -m 'initial commit'
`

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:           "testdata/scripts",
		TestWork:      os.Getenv("TESTWORK") == "1",
		UpdateScripts: os.Getenv("UPDATE_SCRIPTS") == "1",
		Setup: func(env *testscript.Env) error {
			c := exec.Command("/bin/sh", "-c", initRepo)
			c.Dir = env.WorkDir
			out, err := c.CombinedOutput()
			if err != nil {
				env.T().Fatal(string(out), err)
			}
			return nil
		},
	})
}
