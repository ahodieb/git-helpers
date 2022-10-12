package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Git struct {
	WorkingDir string
}

// Rebase the target branch on top of base branch
// `git rebase <base> <target>`
// note that after running `git rebase <base> <target>` the checked out branch becomes <target>
// so we have to check out the current branch again at the end
//
// TODO: handle failures in rebase --abort
func (git *Git) Rebase(base, target string) error {
	currentBranch, err := git.CurrentBranch()
	if err != nil {
		return err
	}

	defer git.Checkout(currentBranch)

	_, stderr, err := git.Exec("rebase", base, target)
	if err != nil {
		git.Exec("rebase", "--abort")

		if stderr != "" {
			return fmt.Errorf("%s, %w", stderr, err)
		}

		return err
	}

	return nil
}

// FindMainBranch returns the name of the main branch
//
// right now it only looks if either `main` or `master` exist
// I can later add more generic code to figure out the default branch if it was neither `master` nor `main`
func (git *Git) FindMainBranch() (string, error) {
	branches, err := git.ListBranches()
	if err != nil {
		return "", err
	}

	for _, b := range branches {
		if b == "main" || b == "master" {
			return b, nil
		}
	}

	return "", fmt.Errorf("did not find `main` or `master` branch, found: %v", branches)
}

// CurrentBranch returns the name of the current checked out branch
// `git branch --show-current`
func (git *Git) CurrentBranch() (string, error) {
	stdout, stderr, err := git.Exec("branch", "--show-current")
	if err != nil {
		return "", err
	}

	if stderr != "" {
		return "", fmt.Errorf(stderr)
	}

	return strings.TrimSpace(stdout), nil
}

// Checkout checks out a <branch>
// `git checkout <branch>`
func (git *Git) Checkout(branch string) error {
	_, stderr, err := git.Exec("checkout", branch)
	if err != nil {
		return err
	}

	if stderr != "" {
		return fmt.Errorf(stderr)
	}

	return nil
}

// ListBranches returns the names of local branches
func (git *Git) ListBranches() ([]string, error) {
	stdout, stderr, err := git.Exec("branch", `--format=%(refname:short)`)
	if err != nil {
		return nil, err
	}

	if stderr != "" {
		return nil, fmt.Errorf(stderr)
	}

	stdout = strings.TrimSpace(stdout)
	return strings.Split(stdout, "\n"), nil
}

// Exec runs a git command with the passed args
func (git *Git) Exec(arg ...string) (string, string, error) {
	var stdout strings.Builder
	var stderr strings.Builder

	arg = append([]string{"-C", git.WorkingDir}, arg...)

	cmd := exec.Command("git", arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}
