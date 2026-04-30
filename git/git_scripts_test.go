package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"git.Exec": func() int {
			out, _ := Exec(os.Args[1:]...)
			fmt.Fprint(os.Stdout, out.Stdout)
			fmt.Fprint(os.Stderr, out.Stderr)
			return out.ExitCode
		},
		"git.CurrentBranch": func() int {
			g := Git{WorkingDir: "."}
			branch, err := g.CurrentBranch()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return exitCode(err)
			}
			fmt.Println(branch)
			return 0
		},
		"git.ListBranches": func() int {
			g := Git{WorkingDir: "."}
			branches, err := g.ListBranches()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return exitCode(err)
			}
			for _, b := range branches {
				fmt.Println(b)
			}
			return 0
		},
	}))
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:           "testdata/scripts",
		TestWork:      os.Getenv("TESTWORK") == "1",
		UpdateScripts: os.Getenv("UPDATE_SCRIPTS") == "1",
	})
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	var exErr *exec.ExitError
	if errors.As(err, &exErr) {
		return exErr.ExitCode()
	}
	return -1
}
