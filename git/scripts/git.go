package scripts

import (
	"path"
	"testing"
)

type Repository struct {
	Dir string
}

type InitOptions struct {
	Name   string
	Branch string
}

func WithBranch(b string) func(opts *InitOptions) error {
	return func(opts *InitOptions) error {
		opts.Branch = b
		return nil
	}
}

func apply[T any](o *T, funcs ...func(*T) error) error {
	for _, f := range funcs {
		if err := f(o); err != nil {
			return err
		}
	}
	return nil
}

// Init initialize a repository and return its path
func Init(t *testing.T, optFunc ...func(*InitOptions) error) Repository {
	t.Helper()

	opts := InitOptions{
		Name:   "testing.repo",
		Branch: "main",
	}

	if err := apply(&opts, optFunc...); err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	ExecScript(t, Params{
		Dir: dir,
		Env: []string{"name=" + opts.Name, "branch=" + opts.Branch},
	}, `
		git init -b "${branch}" "${name}"
		cd "${name}"
		git config user.email 'testing@go.local'
		git config user.name 'git-helper tests'
		git config commit.gpgsign false
		git commit --allow-empty -m '🎉 initial commit'
	`)

	return Repository{Dir: path.Join(dir, "testing.repo")}
}

func (r *Repository) CreateBranch(t *testing.T, branches ...string) {
	t.Helper()
	for _, branch := range branches {
		ExecScript(t, Params{Dir: r.Dir, Env: []string{"branch=" + branch}}, `
		git branch "${branch}"
	`)
	}
}

func (r *Repository) Checkout(t *testing.T, branch string) {
	t.Helper()
	ExecScript(t, Params{Dir: r.Dir, Env: []string{"branch=" + branch}}, `
		git checkout "${branch}"
	`)
}
