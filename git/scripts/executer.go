package scripts

import (
	"os"
	"os/exec"
	"path"
	"testing"
)

type Params struct {
	Dir string
	Env []string
}

func ExecScript(t *testing.T, p Params, script string) string {
	t.Helper()
	dir := t.TempDir()
	scriptPath := path.Join(dir, "script.sh")
	if err := os.WriteFile(scriptPath, []byte(script), 0700); err != nil {
		t.Fatalf("failed to write script to %s, %v", scriptPath, err)
	}

	cmd := exec.Command("/bin/sh", scriptPath)
	cmd.Dir = p.Dir
	cmd.Env = p.Env

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(string(out), err)
	}

	return string(out)
}
