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

// Replay rebases target onto base by writing new refs only — no
// checkout, no working-tree changes. Uses `git replay` (>= 2.44) and
// pipes its output to `git update-ref --stdin`. Returns nil when the
// branch is already up to date with base.
func (g *Git) Replay(base, target string) error {
	out, err := g.Exec("replay", "--onto", base, base+".."+target)
	if err != nil {
		return err
	}
	if out.TrimmedStdout() == "" {
		return nil
	}

	cmd := exec.Command("git", "-C", g.WorkingDir, "update-ref", "--stdin")
	cmd.Stdin = strings.NewReader(out.Stdout)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return fmt.Errorf("update-ref: %s, %w", msg, err)
		}
		return fmt.Errorf("update-ref: %w", err)
	}
	return nil
}

// `git branch --format=%(refname:short)`
func (g *Git) ListBranches() ([]string, error) {
	out, err := g.Exec("branch", `--format=%(refname:short)`)
	if err != nil {
		return nil, err
	}
	return out.StdoutLines(), nil
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
