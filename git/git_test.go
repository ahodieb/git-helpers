package git

import (
	"path/filepath"
	"testing"

	"github.com/ahodieb/git-helpers/git/scripts"
	"github.com/stretchr/testify/assert"
)

// (*Git).Exec must not treat stderr output as an error when the command
// itself succeeded. Git uses stderr as a status channel ("Switched to
// branch...", progress) on exit 0, and earlier code mistakenly surfaced
// that as an error.
func TestExecDoesNotErrorOnStderrWhenCommandSucceeds(t *testing.T) {
	r := scripts.Init(t)
	g := Git{WorkingDir: r.Dir}

	_, err := g.Exec("checkout", "-b", "feature")
	assert.NoError(t, err)
}

// (*Git).Exec must apply WorkingDir via -C. The earlier implementation
// built the exec.Command before prepending -C, so the flag was silently
// dropped and the command ran against the test process's cwd instead.
func TestExecRespectsWorkingDir(t *testing.T) {
	r := scripts.Init(t)
	g := Git{WorkingDir: r.Dir}

	out, err := g.Exec("rev-parse", "--show-toplevel")
	assert.NoError(t, err)

	expected, err := filepath.EvalSymlinks(r.Dir)
	assert.NoError(t, err)
	got, err := filepath.EvalSymlinks(out.TrimmedStdout())
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestListBranches(t *testing.T) {
	r := scripts.Init(t)
	r.CreateBranch(t, "b1", "b2", "b3")
	g := Git{WorkingDir: r.Dir}

	expected := []string{"main", "b1", "b2", "b3"}
	got, err := g.ListBranches()
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, got)
}

