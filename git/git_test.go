package git

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// (*Git).Exec must not treat stderr output as an error when the
// command itself succeeded. Git uses stderr as a status channel
// ("Switched to branch 'foo'", progress) on exit 0; earlier code
// surfaced that as an error.
func TestExecDoesNotErrorOnStderrWhenCommandSucceeds(t *testing.T) {
	g := Git{WorkingDir: setupRepo(t)}
	_, err := g.Exec("checkout", "-b", "feature")
	assert.NoError(t, err)
}

// (*Git).Exec must apply WorkingDir via -C. The earlier implementation
// built exec.Command before prepending -C, so the flag was silently
// dropped and the command ran against the test process's cwd.
func TestExecRespectsWorkingDir(t *testing.T) {
	dir := setupRepo(t)
	g := Git{WorkingDir: dir}

	out, err := g.Exec("rev-parse", "--show-toplevel")
	assert.NoError(t, err)

	expected, err := filepath.EvalSymlinks(dir)
	assert.NoError(t, err)
	got, err := filepath.EvalSymlinks(out)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

// (*Git).Exec must propagate non-zero exit as an error.
func TestExecReturnsErrorOnFailure(t *testing.T) {
	g := Git{WorkingDir: setupRepo(t)}
	_, err := g.Exec("rev-parse", "--verify", "does-not-exist")
	assert.Error(t, err)
}

func TestListBranches(t *testing.T) {
	g := Git{WorkingDir: setupRepo(t, "b1", "b2", "b3")}
	got, err := g.ListBranches()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"main", "b1", "b2", "b3"}, got)
}

// setupRepo creates a real git repo in t.TempDir() with one initial
// commit on `main` and any extra branches forked off it.
func setupRepo(t *testing.T, branches ...string) string {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "repo")

	script := `
		set -e
		git init -q -b main "$1"
		cd "$1"
		git config user.email 'test@git-helpers.local'
		git config user.name 'git-helpers tests'
		git config commit.gpgsign false
		git commit -q --allow-empty -m 'initial commit'
	`
	for _, b := range branches {
		script += "git branch " + b + "\n"
	}

	cmd := exec.Command("/bin/sh", "-c", script, "sh", dir)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("setup: %v\n%s", err, out)
	}
	return dir
}
