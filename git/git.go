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

// Rebase target on top of base. `git rebase <base> <target>` leaves
// HEAD on <target>, so we restore the original branch on the way out.
func (g *Git) Rebase(base, target string) error {
	original, err := g.CurrentBranch()
	if err != nil {
		return err
	}
	defer g.ExecSilent("checkout", original)

	if err := g.ExecSilent("rebase", base, target); err != nil {
		g.ExecSilent("rebase", "--abort")
		return err
	}
	return nil
}

func (g *Git) FindMainBranch() (string, error) {
	branches, err := g.ListBranches()
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

// `git branch --show-current`
func (g *Git) CurrentBranch() (string, error) {
	out, err := g.Exec("branch", "--show-current")
	if err != nil {
		return "", err
	}
	return out.TrimmedStdout(), nil
}

// `git branch --format=%(refname:short)`
func (g *Git) ListBranches() ([]string, error) {
	out, err := g.Exec("branch", `--format=%(refname:short)`)
	if err != nil {
		return nil, err
	}
	return out.StdoutLines(), nil
}

// Cherry returns the +/- prefixed lines from `git cherry -v`
// (commits on HEAD vs its upstream). Empty = no unpublished commits.
func (g *Git) Cherry() ([]string, error) {
	out, err := g.Exec("cherry", "-v")
	if err != nil {
		return nil, err
	}
	return out.StdoutLines(), nil
}

// `git config --global alias.<alias> <command>`
func (g *Git) AddAlias(alias, command string) error {
	key := fmt.Sprintf("alias.%s", alias)
	return g.ExecSilent("config", "--global", key, command)
}

// ExecSilent runs a git command and discards its output.
func (g *Git) ExecSilent(arg ...string) error {
	_, err := g.Exec(arg...)
	return err
}

// Exec runs `git -C <WorkingDir> <args>...` and returns the captured Output.
func (g *Git) Exec(arg ...string) (Output, error) {
	arg = append([]string{"-C", g.WorkingDir}, arg...)
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
	s := o.TrimmedStdout()
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// Exec runs the git binary with the given args. The error is non-nil
// iff git exited non-zero — stderr alone (e.g. "Switched to branch
// 'foo'") does not mark a successful command as failed. On failure the
// error is enriched with stderr for context.
func Exec(arg ...string) (Output, error) {
	var stdout, stderr strings.Builder

	cmd := exec.Command("git", arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()
	out := Output{Stdout: stdout.String(), Stderr: stderr.String()}

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
