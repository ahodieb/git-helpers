package scripts

import (
	"os/exec"
	"path/filepath"
	"testing"
)

type Repository struct {
	Dir string
}

// Init creates a real git repo in t.TempDir() with one initial commit on `main`.
func Init(t *testing.T) Repository {
	return InitOnBranch(t, "main")
}

// InitOnBranch is Init with a custom default-branch name.
func InitOnBranch(t *testing.T, branch string) Repository {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "testing.repo")
	run(t, "", `
		git init -q -b "$1" "$2"
		cd "$2"
		git config user.email 'test@git-helpers.local'
		git config user.name 'git-helpers tests'
		git config commit.gpgsign false
		git commit -q --allow-empty -m 'initial commit'
	`, branch, dir)
	return Repository{Dir: dir}
}

func (r *Repository) CreateBranch(t *testing.T, branches ...string) {
	t.Helper()
	for _, b := range branches {
		run(t, r.Dir, `git branch "$1"`, b)
	}
}

func (r *Repository) Checkout(t *testing.T, branch string) {
	t.Helper()
	run(t, r.Dir, `git checkout "$1"`, branch)
}

// run executes script via /bin/sh -c with positional args ($1, $2, ...).
// Calls t.Fatal on non-zero exit.
func run(t *testing.T, dir, script string, args ...string) {
	t.Helper()
	c := exec.Command("/bin/sh", append([]string{"-c", script, "sh"}, args...)...)
	c.Dir = dir
	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("script failed: %v\n%s", err, out)
	}
}
