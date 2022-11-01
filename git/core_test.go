package git

import (
	"os/exec"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestCore(t *testing.T) {
	_, err := exec.LookPath("git")
	if err != nil {
		t.Fatal("git is not found")
	}

	testscript.Run(t, testscript.Params{
		Dir:         "testdata/scripts/core",
		WorkdirRoot: t.TempDir(),
	})
}
