package git

import (
	"errors"
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
	out, err := git.Exec("branch", "--show-current")
	if err != nil {
		return "", err
	}

	return out.TrimmedStdout(), nil
}

// ListBranches returns the names of local branches
// `git branch --format=%(refname:short)`
func (git *Git) ListBranches() ([]string, error) {
	out, err := git.Exec("branch", `--format=%(refname:short)`)
	if err != nil {
		return nil, err
	}

	return out.StdoutLines(), nil
}

// Cherry returns the diff from branch and upstream
// `git cherry -v`
func (git *Git) Cherry() ([]string, error) {
	out, err := git.Exec("cherry", "-v")
	if err != nil {
		return nil, err
	}

	var changes []string
	for _, line := range out.StdoutLines() {
		if line = strings.TrimSpace(line); line != "" {
			changes = append(changes, line)
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

// ExecSilent runs a git command with the passed args and discards the output.
func (git *Git) ExecSilent(arg ...string) error {
	_, err := git.Exec(arg...)
	return err
}

// Exec runs `git -C <WorkingDir> <args>...` and returns the captured Output.
// The error is non-nil iff git exited non-zero; on failure the message is
// enriched with stderr for context.
func (git *Git) Exec(arg ...string) (Output, error) {
	arg = append([]string{"-C", git.WorkingDir}, arg...)
	return Exec(arg...)
}

type Output struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func (o *Output) TrimmedStdout() string {
	return strings.TrimSpace(o.Stdout)
}

func (o *Output) StdoutLines() []string {
	return strings.Split(o.TrimmedStdout(), "\n")
}

// Exec runs the git binary with the given args and captures its output.
// Errors are signalled by exit code only — git writes status messages to
// stderr on success ("Switched to branch...", progress), so stderr alone
// does not make a successful command an error. On failure the error
// message includes any stderr captured.
func Exec(arg ...string) (Output, error) {
	var stdout strings.Builder
	var stderr strings.Builder

	cmd := exec.Command("git", arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	out := Output{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if runErr == nil {
		return out, nil
	}

	var exErr *exec.ExitError
	if errors.As(runErr, &exErr) {
		out.ExitCode = exErr.ExitCode()
	} else {
		out.ExitCode = 1
	}

	if msg := strings.TrimSpace(out.Stderr); msg != "" {
		return out, fmt.Errorf("%s, %w", msg, runErr)
	}
	return out, runErr
}
