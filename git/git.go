package git

import (
	"os/exec"
	"strings"
)

type Git struct {
	WorkingDir string
}

func (git *Git) Exec(arg ...string) (string, string, error) {
	var stdout strings.Builder
	var stderr strings.Builder

	arg = append([]string{"-C", git.WorkingDir}, arg...)

	cmd := exec.Command("git", arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
