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
			fmt.Println(branch)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return exitCode(err)
		},
		"git.ListBranches": func() int {
			g := Git{WorkingDir: "."}
			branches, err := g.ListBranches()
			for _, b := range branches {
				fmt.Println(b)
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			return exitCode(err)
		},
	}))
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:           "testdata/scripts",
		TestWork:      true,
		UpdateScripts: true,
		//Cmds: map[string]func(ts *testscript.TestScript, neg bool, args []string){
		//	"gitt":,
		//},
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
