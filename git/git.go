package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Git is a thin wrapper around the git CLI, bound to a working directory.
type Git struct {
	WorkingDir string
}

// Exec runs `git -C <WorkingDir> <args>...` and returns trimmed stdout.
// Errors are signalled by exit code only — git writes status messages
// to stderr on success, so stderr alone is not a failure signal. On
// failure the error is enriched with stderr.
func (g *Git) Exec(arg ...string) (string, error) {
	out, err := g.run("", arg...)
	return strings.TrimSpace(out), err
}

// ListBranches returns local branch names.
func (g *Git) ListBranches() ([]string, error) {
	out, err := g.Exec("branch", `--format=%(refname:short)`)
	if err != nil || out == "" {
		return nil, err
	}
	return strings.Split(out, "\n"), nil
}

// Replay rebases target onto base by writing new refs only — no
// checkout, no working-tree changes. Uses `git replay` (>= 2.44) and
// pipes its raw output to `git update-ref --stdin` (which is
// line-terminated, so the trailing newline must survive).
func (g *Git) Replay(base, target string) error {
	refs, err := g.run("", "replay", "--onto", base, base+".."+target)
	if err != nil {
		return err
	}
	if strings.TrimSpace(refs) == "" {
		return nil
	}
	if _, err := g.run(refs, "update-ref", "--stdin"); err != nil {
		return fmt.Errorf("update-ref: %w", err)
	}
	return nil
}

func (g *Git) run(stdin string, arg ...string) (string, error) {
	var stdout, stderr strings.Builder
	cmd := exec.Command("git", append([]string{"-C", g.WorkingDir}, arg...)...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			return "", fmt.Errorf("%s, %w", msg, err)
		}
		return "", err
	}
	return stdout.String(), nil
}
