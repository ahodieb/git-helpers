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

	err = git.ExecSilent("rebase", base, target)
	if err != nil {
		git.Exec("rebase", "--abort")
	}

	return err
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
	stdout, err := git.Exec("branch", "--show-current")
	if err != nil {
		return "", err
	}

	return stdout, nil
}

// ListBranches returns the names of local branches
// `git branch --format=%(refname:short)
func (git *Git) ListBranches() ([]string, error) {
	stdout, err := git.Exec("branch", `--format=%(refname:short)`)
	if err != nil {
		return nil, err
	}

	return strings.Split(stdout, "\n"), nil
}

// Cherry returns the diff from branch and upstream
// `git cherry -v`
func (git *Git) Cherry() ([]string, error) {
	stdout, err := git.Exec("cherry", "-v")
	if err != nil {
		return nil, err
	}

	var changes []string
	for _, change := range strings.Split(stdout, "\n") {
		if strings.TrimSpace(change) != "" {
			changes = append(changes, change)
		}
	}
	return changes, nil
}

// Checkout checks out a <branch>
// `git checkout <branch>`
func (git *Git) Checkout(branch string) error {
	return git.ExecSilent("checkout", branch)
}

// AddAlias adds a global alias to git config
// `git config --global alias.<alias> 'command'
func (git *Git) AddAlias(alias, command string) error {
	alias = fmt.Sprintf("alias.%s", alias)
	return git.ExecSilent("config", "--global", alias, command)
}

// ExecSilent runs a git command with the passed args and ignores the output
func (git *Git) ExecSilent(arg ...string) error {
	_, err := git.Exec(arg...)
	return err
}

// Exec runs a git command with the passed args
func (git *Git) Exec(arg ...string) (string, error) {
	arg = append([]string{"-C", git.WorkingDir}, arg...)
	cmd := exec.Command("git", arg...)

	var stdout strings.Builder
	var stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.String() != "" {
			return "", fmt.Errorf("%s, %w", stderr.String(), err)
		}

		if stdout.String() != "" {
			return "", fmt.Errorf("%s, %w", stdout.String(), err)
		}

		return "", err
	}

	//if stderr.String() != "" {
	//	_, _ = fmt.Fprintln(os.Stderr, "git:", stderr.String())
	//}

	return strings.TrimSpace(stdout.String()), nil
}
